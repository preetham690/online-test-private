package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"onlinetest/database"
	"onlinetest/normaluser/models"

	"gopkg.in/mgo.v2/bson"
)

func GetTestByCatHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the category from the URL parameters
	category := r.URL.Query().Get("category")

	if category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	// Fetch the user test details
	userTestDetails, err := GetTestByCat("mongodb://localhost:27017", "testDB", "Test_Details", category)
	if err != nil {
		http.Error(w, "Error fetching user test details", http.StatusInternalServerError)
		return
	}

	// Return the user test details as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userTestDetails); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetTestByCat(dbURI, dbName, collectionName, category string) ([]models.UserTestDetails, error) {
	// Initialize the MongoDB client and collection
	client, collection, err := database.InitMongoDB(dbURI, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Define a filter to find the tests by category
	filter := bson.M{"category": category}

	// Find all test documents in the collection that match the filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var userTestDetails []models.UserTestDetails

	// Iterate through the cursor to fetch all matching test documents
	for cursor.Next(context.Background()) {
		var test models.UserTestDetails
		if err := cursor.Decode(&test); err != nil {
			return nil, err
		}
		userTestDetails = append(userTestDetails, test)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userTestDetails, nil
}
