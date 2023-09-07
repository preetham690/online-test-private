package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func InitMongoDB() (*mongo.Client, *mongo.Collection, error) {
// 	// Load sensitive information from environment variables
// 	os.Setenv("DB_URI", "mongodb://localhost:27017")
// 	os.Setenv("DB_NAME", "auth")

// 	dbURI := os.Getenv("DB_URI")
// 	dbName := os.Getenv("DB_NAME")
// 	if dbURI == "" || dbName == "" {
// 		return nil, nil, fmt.Errorf("DB_URI and DB_NAME environment variables are required")
// 	}

// 	// Connect to MongoDB
// 	opts := options.Client().ApplyURI(dbURI)
// 	client, err := mongo.Connect(context.Background(), opts)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("error connecting to MongoDB: %v", err)
// 	}

// 	userCollection := client.Database(dbName).Collection("users")
// 	if userCollection == nil {
// 		return nil, nil, fmt.Errorf("failed to initialize user collection")
// 	}

//		return client, userCollection, nil
//	}
func InitMongoDB(dbURI, dbName, collectionName string) (*mongo.Client, *mongo.Collection, error) {
	if dbURI == "" || dbName == "" || collectionName == "" {
		return nil, nil, fmt.Errorf("DB_URI, DB_NAME, and CollectionName must not be empty")
	}

	// Connect to MongoDB
	opts := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Access the specified collection
	collection := client.Database(dbName).Collection(collectionName)
	if collection == nil {
		return nil, nil, fmt.Errorf("failed to initialize collection: %s", collectionName)
	}

	return client, collection, nil
}
