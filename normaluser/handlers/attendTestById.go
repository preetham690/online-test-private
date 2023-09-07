package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"onlinetest/database"
	"onlinetest/normaluser/models"
	"strconv"
)

func SubmitUserAnswersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the email parameter from the URL
	email := r.URL.Query().Get("email")

	// Extract the testID parameter from the URL
	testIDstr := r.URL.Query().Get("testID")
	if testIDstr == "" {
		http.Error(w, "Test id is required", http.StatusBadRequest)
		return
	}

	// Parse the testID from the URL
	testID, err := strconv.Atoi(testIDstr)
	if err != nil {
		http.Error(w, "Test id is required", http.StatusBadRequest)
		return
	}

	// Parse the user's answers from the request body
	var userAnswers models.UserAnswer
	if err := json.NewDecoder(r.Body).Decode(&userAnswers); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Validate that the testID in the userAnswers matches the testID from the URL
	if userAnswers.TestID != testID {
		http.Error(w, "Attending the wrong test", http.StatusBadRequest)
		return
	}

	// Set the email in the UserAnswer struct
	userAnswers.Email = email

	// Store the user's answers in the MongoDB database
	if err := StoreUserAnswers("mongodb://localhost:27017", "testDB", "UserAnswers", userAnswers); err != nil {
		http.Error(w, "Error storing user answers", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User answers stored successfully")
}

func StoreUserAnswers(dbURI, dbName, collectionName string, userAnswers models.UserAnswer) error {
	// Initialize the MongoDB client and collection
	client, collection, err := database.InitMongoDB(dbURI, dbName, collectionName)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Insert the user's answers into the collection
	_, err = collection.InsertOne(context.Background(), userAnswers)
	if err != nil {
		return err
	}

	return nil
}
