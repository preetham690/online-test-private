package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"onlinetest/auth/handlers"

	"github.com/gorilla/mux"
)

// var (
// 	PORT string
// 	//logger *log.Logger
// )

// func main() {

// 	// flag.StringVar(&PORT, "port", "20200", "--port=20200 or -port=20200")
// 	// flag.Parse()

// 	// router := mux.NewRouter()
// 	// srv := http.Server

// 	// HTTP server setup
// 	http.HandleFunc("/register", hand.Register)
// 	http.HandleFunc("/login", hand.Login)

// 	port := 20200
// 	fmt.Printf("Server is running on port %d\n", port)
// 	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

// 	if err != nil {
// 		log.Fatalf("Error starting HTTP server: %v", err)
// 	}
// }

func main() {
	//port
	var port int
	flag.IntVar(&port, "port", 20200, "Port number for the HTTP server")
	flag.Parse()

	//gorilla mux
	router := mux.NewRouter()
	//calling http handlers using mux
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	fmt.Printf("Server is running on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
