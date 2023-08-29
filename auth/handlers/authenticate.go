package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"onlinetest/auth/models"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// var client *mongo.Client
// var userCollection *mongo.Collection

func initMongoDB() (*mongo.Client, *mongo.Collection, error) {
	// Load sensitive information from environment variables
	os.Setenv("DB_URI", "mongodb+srv://VictoriaSecretsPreS:columbus@cluster0.yikox0z.mongodb.net/?retryWrites=true&w=majority")
	os.Setenv("DB_NAME", "<dbname>")

	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	if dbURI == "" || dbName == "" {
		return nil, nil, fmt.Errorf("DB_URI and DB_NAME environment variables are required")
	}

	// Connect to MongoDB
	opts := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	userCollection := client.Database(dbName).Collection("users")
	if userCollection == nil {
		return nil, nil, fmt.Errorf("failed to initialize user collection")
	}

	return client, userCollection, nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	// // Load sensitive information from environment variables
	// os.Setenv("DB_URI", "mongodb+srv://VictoriaSecretsPreS:columbus@cluster0.yikox0z.mongodb.net/?retryWrites=true&w=majority")
	// os.Setenv("DB_NAME", "<dbname>")

	// dbURI := os.Getenv("DB_URI")
	// dbName := os.Getenv("DB_NAME")
	// if dbURI == "" || dbName == "" {
	// 	log.Fatal("DB_URI and DB_NAME environment variables are required")
	// }

	// // Connect to MongoDB
	// opts := options.Client().ApplyURI(dbURI)
	// client, err := mongo.Connect(context.Background(), opts)
	// if err != nil {
	// 	log.Fatalf("Error connecting to MongoDB: %v", err)
	// }
	// defer client.Disconnect(context.Background())

	// userCollection = client.Database(dbName).Collection("users")
	// if userCollection == nil {
	// 	log.Fatal("Failed to initialize user collections")
	// }

	//calling the mongodb init function
	client, userCollection, err := initMongoDB()
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
	client, userCollection, err := initMongoDB()
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

	// Check if the user with the provided email exists in the database
	u := new(models.User)
	err = userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(u)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found in the database, return an appropriate message
			log.Printf("Error finding user: %v", err)
			http.Error(w, "User not found. Please register first.", http.StatusUnauthorized)
			return
		} else {
			// Some other error occurred
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Compare the hashed password with the password entered
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// Passwords don't match, return an appropriate message
		http.Error(w, "Invalid Username Or Password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Authentication Successful")
}
