package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error while loading env file")
	}

	MongoDBURL := os.Getenv("MONGODB_URL")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MongoDBURL).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// Send a ping to confirm a successful connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error while pinging MongoDB")
	}

	fmt.Println("✅ Successfully connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("redxcoder").Collection(collectionName)

	// Test query to check if the collection is accessible
	_, err := collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		log.Fatalf("Failed to access collection %s: %v", collectionName, err)
	}

	return collection
}
