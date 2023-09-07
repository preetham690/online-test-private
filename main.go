package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	admin "onlinetest/admin/handlers"
	auth "onlinetest/auth/handlers"
	"onlinetest/middleware"
	normal "onlinetest/normaluser/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Port
	var port int
	flag.IntVar(&port, "port", 20201, "Port number for the HTTP server")
	flag.Parse()

	// Gorilla Mux
	router := mux.NewRouter()

	// Calling HTTP handlers using Mux
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", auth.Register).Methods("POST")
	authRouter.HandleFunc("/login", auth.Login).Methods("POST")

	// we are making use of subrouter so that we can pass the middleware only for admin handlers
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.CheckUserRoleMiddleware)
	adminRouter.HandleFunc("/create-test", admin.CreateTest).Methods("POST")
	adminRouter.HandleFunc("/create-cat", admin.CreateCategory).Methods("POST")
	adminRouter.HandleFunc("/create-dashboard", admin.CreateDashboard).Methods("POST")

	//and this for normaluser folder
	normalRouter := router.PathPrefix("/normaluser").Subrouter()
	normalRouter.HandleFunc("/get-cat", normal.GetCatHandler).Methods("GET")
	normalRouter.HandleFunc("/get-cat-by-id", normal.GetCatByIdHandler).Methods("GET")
	normalRouter.HandleFunc("/get-test-by-cat", normal.GetTestByCatHandler).Methods("GET")
	normalRouter.HandleFunc("/get-test-by-id", normal.GetTestByIdHandler).Methods("GET")
	normalRouter.HandleFunc("/attend-test-by-id", normal.SubmitUserAnswersHandler).Methods("POST")
	normalRouter.HandleFunc("/get-result", normal.CalculateAndStoreUserResult).Methods("GET")

	fmt.Printf("Server is running on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
