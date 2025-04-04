package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dayalubajpai/jwtlearninggo/database"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	User_type  string
	UID        string
	jwt.StandardClaims
}

var collection *mongo.Collection = database.OpenCollection(database.Client, "redxcoder")

var SECRET_KEY string

func init() {
	SECRET_KEY = os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		log.Fatal("Secret key not set in environment variables")
	}
}

func GenerateAllTokens(email string, firstName string, lastName string, userType string, userId string) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		UID:        userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			Issuer:    "boracay",
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(48)).Unix(),
			Issuer:    "boracay",
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access Token: %w", err)
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", fmt.Errorf("failed to generate access refreshToken: %w", err)
	}

	return token, refreshToken, nil
}

func ValidateToken(token string) (claims *SignedDetails, msg string) {

	tokenClaims, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	var ok bool
	if claims, ok = tokenClaims.Claims.(*SignedDetails); !ok {
		msg = "Token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(token string, refreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    time.Now(),
		},
	}
	upsert := true
	filter := bson.M{"user_id": userId}

	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := collection.UpdateOne(ctx, filter, update, opt)

	if err != nil {
		log.Panic(err)
		return
	} else {
		log.Println("User updated successfully")
		return
	}

}
