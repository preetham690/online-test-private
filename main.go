package main

import (
	"fmt"
	"log"
	"net/http"
	hand "onlinetest/auth/handlers"
)

func main() {

	// HTTP server setup
	http.HandleFunc("/register", hand.Register)
	http.HandleFunc("/login", hand.Login)

	port := 20200
	fmt.Printf("Server is running on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
