package main

import (
	"log"
	"net/http"

	"github.com/ombima56/insights-edge/contracts"
	"github.com/ombima56/insights-edge/database"
	"github.com/ombima56/insights-edge/handlers"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// Initialize contract handler
	_ = contracts.NewSmartContractHandler(db)

	// Setup static file server
	handlers.SetupStaticFileServer()

	// Register HTTP handlers
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/api/signup", handlers.SignupHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/profile", handlers.GetUserProfileHandler)
	http.HandleFunc("/api/users", handlers.GetUsersHandler)

	// Start the server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
