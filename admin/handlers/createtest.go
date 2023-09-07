package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"onlinetest/admin/models"
	"onlinetest/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateTestHandler handles the creation of a new test.
func CreateTest(w http.ResponseWriter, r *http.Request) {
	var newTest models.TestDetails

	// Decode the request body to get test data
	if err := json.NewDecoder(r.Body).Decode(&newTest); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Generate a random test ID
	newTest.TestID = generateUniqueTestID()

	// Initializing the MongoDB client and collection
	client, userCollection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Details")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Check if the selected category exists or is valid
	categoryName := newTest.Category
	categoryExists, err := ValidateDashboard(categoryName)
	if err != nil {
		http.Error(w, "Error validating category", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Selected category does not exist or is invalid", http.StatusBadRequest)
		return
	}

	// Get dashboard details by category
	dashboard, err := GetDashboardByCategory(categoryName)
	if err != nil {
		http.Error(w, "Error getting dashboard data", http.StatusInternalServerError)
		return
	}

	if len(newTest.Questions) != dashboard.NumberOfQuestions {
		http.Error(w, "Number of questions in the test does not match the dashboard", http.StatusBadRequest)
		return
	}

	// Inserting test details into the a new collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, newTest)
	if err != nil {
		http.Error(w, "Error inserting test data", http.StatusInternalServerError)
		return
	}

	// Return a success response
	response := map[string]interface{}{
		"message":      "Test created successfully",
		"insertedData": result,
	}

	//here we just setting the header nothing else
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// generates a random number between 0-1000
func generateUniqueTestID() int {
	return rand.Intn(1000)
}

// for validation
func ValidateDashboard(categoryName string) (bool, error) {
	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Dashboard")
	if err != nil {
		return false, err
	}
	defer client.Disconnect(context.Background())

	filter := bson.M{"category": categoryName}
	count, err := collection.CountDocuments(context.Background(), filter, nil)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func GetDashboardByCategory(categoryName string) (*models.Dashboard, error) {
	client, userCollection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Dashboard")
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	filter := bson.M{"category": categoryName}
	var dashboard models.Dashboard
	err = userCollection.FindOne(context.Background(), filter).Decode(&dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}
