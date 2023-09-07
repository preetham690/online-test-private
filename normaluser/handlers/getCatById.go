package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"onlinetest/admin/models"
	"onlinetest/database"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

func GetCatByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the category ID from the request (e.g., from URL parameters or request body)
	categoryID := r.URL.Query().Get("id")

	if categoryID == "" {
		http.Error(w, "Category ID is missing", http.StatusBadRequest)
		return
	}

	category, err := GetCategoryByID("mongodb://localhost:27017", "testDB", "Test_Category", categoryID)
	if err != nil {
		http.Error(w, "Error fetching category", http.StatusInternalServerError)
		return
	}

	// Return the category as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetCategoryByID(dbURI, dbName, collectionName string, categoryID string) (*models.TestCategory, error) {
	// Initialize the MongoDB client and collection
	client, collection, err := database.InitMongoDB(dbURI, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Parse the categoryID as an integer
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return nil, err
	}

	// Define a filter to find the category by ID
	filter := bson.M{"id": id}

	// Find the category document in the collection
	var category models.TestCategory
	err = collection.FindOne(context.Background(), filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
