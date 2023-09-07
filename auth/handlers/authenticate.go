package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"onlinetest/auth/models"
	"onlinetest/database"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// var client *mongo.Client
// var userCollection *mongo.Collection

func Register(w http.ResponseWriter, r *http.Request) {

	//calling the mongodb init function
	client, userCollection, err := database.InitMongoDB("mongodb://localhost:27017", "auth", "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	//if otherthan post req it will throw an error
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("unimplemented method"))
		return
	}

	u := new(models.User)

	err = json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	err = u.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user with the provided email already exists in the database
	// existingUser := models.User{}
	// err = userCollection.FindOne(context.Background(), bson.M{"email": u.Email}).Decode(&existingUser)
	// if err == nil {
	// 	// User with the same email already exists, return an error
	// 	http.Error(w, "User with this email already exists", http.StatusConflict)
	// 	return
	// } else if err != mongo.ErrNoDocuments {
	// 	// Some other error occurred
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }
	// err = userCollection.FindOne(context.Background(), bson.M{"mobile": u.Mobile}).Decode(&existingUser)
	// if err == nil {
	// 	// User with the same email already exists, return an error
	// 	http.Error(w, "User with this Mobile already exists", http.StatusConflict)
	// 	return
	// } else if err != mongo.ErrNoDocuments {
	// 	// Some other error occurred
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }
	// Create a WaitGroup to wait for Goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to collect errors from Goroutines
	errCh := make(chan error, 2)

	// Use Goroutines to perform MongoDB operations concurrently
	wg.Add(2)
	go func() {
		defer wg.Done()

		// Check if the user with the provided email already exists
		err := userCollection.FindOne(context.Background(), bson.M{"email": u.Email}).Err()
		if err == nil {
			// User with the same email already exists
			errCh <- fmt.Errorf("user with this email already exists")
		} else if err != mongo.ErrNoDocuments {
			// Some other error occurred
			errCh <- fmt.Errorf("internal server error")
		}
	}()

	go func() {
		defer wg.Done()

		// Check if the user with the provided mobile already exists
		err := userCollection.FindOne(context.Background(), bson.M{"mobile": u.Mobile}).Err()
		if err == nil {
			// User with the same mobile already exists
			errCh <- fmt.Errorf("user with this Mobile already exists")
		} else if err != mongo.ErrNoDocuments {
			// Some other error occurred
			errCh <- fmt.Errorf("internal server error")
		}
	}()

	// Waiting for Goroutines to finish
	wg.Wait()

	// Closing the error channel so that it indicates that all Goroutines are done
	close(errCh)

	// Check for errors from Goroutines
	for err := range errCh {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Hashing the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		User_ID:  primitive.NewObjectID(),
		Name:     u.Name,
		Email:    u.Email,
		Mobile:   u.Mobile,
		Password: string(hashedPass),
		IsAdmin:  u.IsAdmin,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Error inserting user data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
	fmt.Fprintln(w, "User Registered Successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Calling the MongoDB init function
	client, userCollection, err := database.InitMongoDB("mongodb://localhost:27017", "auth", "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// If other than POST req, throw an error
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("unimplemented method"))
		return
	}

	// email := r.FormValue("email")
	// password := r.FormValue("password")

	formData := new(models.LoginForm)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	email := formData.Email
	password := formData.Password

	// // Check if the user with the provided email exists in the database
	// u := new(models.User)
	// err = userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(u)

	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		// User not found in the database, return an appropriate message
	// 		log.Printf("Error finding user: %v", err)
	// 		http.Error(w, "User not found. Please register first.", http.StatusUnauthorized)
	// 		return
	// 	} else {
	// 		// Some other error occurred
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	// // Compare the hashed password with the password entered
	// err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	// if err != nil {
	// 	// Passwords don't match, return an appropriate message
	// 	http.Error(w, "Invalid Username Or Password", http.StatusUnauthorized)
	// 	return
	// }

	// if u.IsAdmin {
	// 	// Admin functionality
	// 	fmt.Fprintln(w, "Admin Authentication Successful")
	// 	// You can add admin-specific logic here
	// } else {
	// 	// Normal user functionality
	// 	fmt.Fprintln(w, "User Authentication Successful")
	// 	// You can add user-specific logic here
	// }
	// Use Goroutines to perform MongoDB operations concurrently
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Check if the user with the provided email exists in the database
		u := new(models.User)
		err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(u)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				// User not found in the database
				log.Printf("Error finding user: %v", err)
				http.Error(w, "User not found. Please register first.", http.StatusUnauthorized)
				return
			} else {
				// Some other error occurred
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		// Comparing the hashed password with the password entered
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			// Passwords don't match
			http.Error(w, "Invalid Username Or Password", http.StatusUnauthorized)
			return
		}

		if u.IsAdmin {
			// Admin access
			fmt.Fprintln(w, "Admin Authentication Successful")
		} else {
			// Normal user
			fmt.Fprintln(w, "User Authentication Successful")
		}
	}()

	// Waiting for Goroutines to finish
	wg.Wait()
}
