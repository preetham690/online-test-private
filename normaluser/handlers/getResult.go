package handlers

import (
	"context"
	"fmt"
	"net/http"
	admin "onlinetest/admin/models"
	"onlinetest/database"
	normal "onlinetest/normaluser/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserAnswersByEmailAndTestID(email string, testID int) (*normal.UserAnswer, error) {
	// Initialize the MongoDB client and collection for UserAnswer
	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "UserAnswers")
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Define a filter to find user answers by email and testID
	filter := bson.M{"email": email, "testID": testID}

	// Find the user answers document in the collection
	var userAnswers normal.UserAnswer
	err = collection.FindOne(context.Background(), filter).Decode(&userAnswers)
	if err != nil {
		return nil, err
	}

	return &userAnswers, nil
}

func GetActualAnswersByTestID(testID int) (*admin.TestDetails, error) {
	// Initialize the MongoDB client and collection for TestDetails
	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "Test_Details")
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Define a filter to find actual answers by testID
	filter := bson.M{"testID": testID}

	// Find the actual answers document in the collection
	var actualAnswers admin.TestDetails
	err = collection.FindOne(context.Background(), filter).Decode(&actualAnswers)
	if err != nil {
		return nil, err
	}

	return &actualAnswers, nil
}

func CalculateUserScore(userAnswers *normal.UserAnswer, actualAnswers *admin.TestDetails) int {
	// Implement code to compare user answers with actual answers
	// Calculate and return the user's score
	score := 0

	// Compare user's answers with actual answers
	for _, userQA := range userAnswers.QuestionAnswers {
		for _, actualQA := range actualAnswers.Questions {
			if userQA.QuestionID == actualQA.QuestionID && userQA.Answer == actualQA.Answer {
				// Increment the score for a correct answer
				score++
				break
			}
		}
	}

	return score
}

func CalculateAndStoreUserResult(w http.ResponseWriter, r *http.Request) {
	// Extract email and testID from request parameters
	email := r.URL.Query().Get("email")
	testID := r.URL.Query().Get("testID")

	if email == "" || testID == "" {
		http.Error(w, "Email and testID are required", http.StatusBadRequest)
		return
	}

	// Convert testID to an integer
	testIDInt := 0
	_, err := fmt.Sscanf(testID, "%d", &testIDInt)
	if err != nil {
		http.Error(w, "Invalid testID", http.StatusBadRequest)
		return
	}

	// Fetch user answers and actual answers
	userAnswers, err := GetUserAnswersByEmailAndTestID(email, testIDInt)
	if err != nil {
		http.Error(w, "Error fetching user answers", http.StatusInternalServerError)
		return
	}

	actualAnswers, err := GetActualAnswersByTestID(testIDInt)
	if err != nil {
		http.Error(w, "Error fetching actual answers", http.StatusInternalServerError)
		return
	}

	// Calculate the user's score
	score := CalculateUserScore(userAnswers, actualAnswers)

	// Store the user result in the UserResult collection
	client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "testDB", "UserResult")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Create a UserResult struct with email, testID, and score
	userResult := normal.UserResult{
		Email:  email,
		TestID: testIDInt,
		Score:  score,
	}

	// Insert the user result into the UserResult collection
	_, err = collection.InsertOne(context.Background(), userResult)
	if err != nil {
		http.Error(w, "Error inserting user result", http.StatusInternalServerError)
		return
	}

	// Return the user's score as a JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"score": %d}`, score)
}
