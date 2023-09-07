package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"onlinetest/admin/models"
	"onlinetest/database"

	"gopkg.in/mgo.v2/bson"
)

// this handler is used get the category from the mongodb
func GetCatHandler(w http.ResponseWriter, r *http.Request) {
	//
	categories, err := GetCategories("mongodb://localhost:27017", "testDB", "Test_Category")
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	// Returns the categories as a JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
func GetCategories(dbURI, dbName, collectionName string) ([]models.TestCategory, error) {
	//here we are calling the database package
	client, collection, err := database.InitMongoDB(dbURI, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	var categories []models.TestCategory

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var category models.TestCategory
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
