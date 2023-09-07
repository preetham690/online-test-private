package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"onlinetest/admin/models"
	"onlinetest/database"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTestByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the test ID from the URL parameters
	testIDstr := r.URL.Query().Get("testID")

	if testIDstr == "" {
		http.Error(w, "Test ID is required", http.StatusBadRequest)
		return
	}

	//parse into string
	testID, err := strconv.Atoi(testIDstr)
	if err != nil {
		http.Error(w, "Invalid Test Id", http.StatusBadRequest)
		return
	}

	// Fetch the user test details by ID
	userTestDetails, err := GetTestByIdDetails("mongodb://localhost:27017", "testDB", "Test_Details", testID)
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

func GetTestByIdDetails(dbURI, dbName, collectionName string, testID int) (*models.TestDetails, error) {
	// Initialize the MongoDB client and collection
	client, collection, err := database.InitMongoDB(dbURI, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Define a filter to find the test by its ID
	filter := bson.M{"testID": testID}

	// Define the projection to exclude the "Answer" field
	projection := bson.M{"questions.answer": 0}

	// Create FindOneOptions with the projection
	opts := options.FindOne().SetProjection(projection)

	// Find the test document in the collection with the specified filter and projection
	var test models.TestDetails
	err = collection.FindOne(context.Background(), filter, opts).Decode(&test)
	if err != nil {
		return nil, err
	}

	return &test, nil
}
