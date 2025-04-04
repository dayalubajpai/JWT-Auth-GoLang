package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dayalubajpai/jwtlearninggo/database"
	helpers "github.com/dayalubajpai/jwtlearninggo/helpers"
	"github.com/dayalubajpai/jwtlearninggo/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection = database.OpenCollection(database.Client, "redxcoder")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("Error while ...:", err)
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		log.Println("Error while comparing password:", err)
		return false, "Invalid password"
	}
	return true, ""
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check email and phone in parallel
		var wg sync.WaitGroup
		var emailCount, phoneCount int64
		var emailErr, phoneErr error

		wg.Add(2)
		go func() {
			defer wg.Done()
			emailCount, emailErr = collection.CountDocuments(ctx, bson.M{"email": user.Email})
		}()
		go func() {
			defer wg.Done()
			phoneCount, phoneErr = collection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		}()
		wg.Wait()

		if emailErr != nil || phoneErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking email or phone"})
			return
		}

		if emailCount > 0 || phoneCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone or Email already exists"})
			return
		}

		// Hash password (not in goroutine since it's CPU-intensive)
		password := HashPassword(*user.Password)
		user.Password = &password

		// Set user metadata
		now := time.Now().Format(time.RFC3339)
		user.Created_at, _ = time.Parse(time.RFC3339, now)
		user.Updated_at, _ = time.Parse(time.RFC3339, now)
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Insert user first
		result, err := collection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item not created"})
			return
		}

		// Generate tokens after successful user creation
		token, refreshToken, _ := helpers.GenerateAllTokens(
			*user.Email,
			*user.First_name,
			*user.Last_name,
			*user.User_type,
			user.User_id,
		)

		// Update user with tokens
		update := bson.M{
			"$set": bson.M{
				"token":         token,
				"refresh_token": refreshToken,
			},
		}

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"user_id": user.User_id},
			update,
		)

		if err != nil {
			// Don't return error since user is created successfully
			log.Printf("Error updating tokens: %v", err)
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User type not matched"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		pageNumber, err := strconv.Atoi(c.Query("pageNumber"))
		if err != nil || pageNumber < 1 {
			pageNumber = 1
		}
		startIndex := (pageNumber - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, // Count the total number of documents
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "users_items", Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		result, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching users"})
			return
		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding users"})
			return
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = collection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, foundUser)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		err := helpers.MatchUserTypeToUID(c, userId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User type not matched"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancel()

		c.JSON(http.StatusOK, user)
	}
}
