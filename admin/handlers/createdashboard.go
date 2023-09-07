package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"onlinetest/admin/models"

	//"onlinetest/Admin/models"
	"onlinetest/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Category")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("unimplemented method"))
		return
	}

	category := models.TestCategory{}
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Insert the category into MongoDB
	_, err = collection.InsertOne(context.Background(), category)
	if err != nil {
		http.Error(w, "Error inserting category data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Category Created Successfully")
}

func CreateDashboard(w http.ResponseWriter, r *http.Request) {

	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Dashboard")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// POST method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("unimplemented method"))
		return
	}

	// Decoding and inserting a single dashboard from request body
	dashboardDetails := models.Dashboard{}
	err = json.NewDecoder(r.Body).Decode(&dashboardDetails)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Check if the selected category exists or is valid
	catName := dashboardDetails.Category
	categoryExists, err := ValidateCategory(catName)
	if err != nil {
		http.Error(w, "Error validating category", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Selected category does not exist or is invalid", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, dashboardDetails)
	if err != nil {
		http.Error(w, "Error inserting dashboard data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
	fmt.Fprintln(w, "Dashboard Created Successfully")

}

// validation function
func ValidateCategory(categoryName string) (bool, error) {
	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Category")
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
