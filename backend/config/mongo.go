package config

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var MongoCtx = context.Background()

// Initialize MongoDB
func InitMongoDB() {
	// Get MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(MongoCtx, clientOptions)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	Client = client

	// Check the connection
	err = client.Ping(MongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}
}
