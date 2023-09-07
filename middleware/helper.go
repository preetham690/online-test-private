package middleware

import (
	"context"
	"fmt"
	"net/http"
	u "onlinetest/auth/models"
	"onlinetest/database"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// This is the middleware
func CheckUserRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//calling mongodb fucntion
		client, collection, err := database.InitMongoDB("mongodb://localhost:27017", "auth", "users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Disconnect(context.Background())

		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email parameter is missing", http.StatusBadRequest)
			return
		}

		// here we are fetching the user details using email
		user, err := GetUserByEmail(email, collection)
		if err != nil {
			http.Error(w, "Error fetching user details", http.StatusInternalServerError)
			return
		}

		// if user != nil && user.IsAdmin {
		// 	w.Header().Set("X-UserRole", "Admin")
		// } else {
		// 	w.Header().Set("X-UserRole", "Normal")
		// }
		var message string

		if user != nil && user.IsAdmin {
			message = "Admin user"
		} else {
			message = "Normal user"
		}

		// Set the response content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Write the message to the response body as JSON
		responseJSON := fmt.Sprintf(`{"message": "%s"}`, message)
		w.Write([]byte(responseJSON))

		// Call the next handler in the chain
		if user != nil && user.IsAdmin {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Access denied. Only admin users can access this endpoint", http.StatusForbidden)
			return
		}
	})
}

func GetUserByEmail(email string, collection *mongo.Collection) (*u.User, error) {
	var user u.User

	// Define a filter to find the user by email
	filter := bson.M{"email": email}

	// Find the user document in the collection
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User with the specified email was not found
			return nil, nil
		}
		// An error occurred while fetching the user
		return nil, err
	}

	// User found
	return &user, nil
}
