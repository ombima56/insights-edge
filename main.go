package main

import (
	"insights-edge/internal/database"
	"insights-edge/internal/routes"
	"insights-edge/internal/util"
	"log"
	"net/http"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	routes.InitRoutes()

	port, err := util.ValidatePort()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server starting on %v\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
