package main

import (
	"log"
	"net/http"

	"github.com/ombima56/insights-edge/internal/database"
	"github.com/ombima56/insights-edge/internal/routes"
	"github.com/ombima56/insights-edge/internal/util"
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

	log.Printf("Server starting on http://localhost%v\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
