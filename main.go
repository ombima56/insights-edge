package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ombima56/insights-edge/internal/database"
	"github.com/ombima56/insights-edge/internal/routes"
	"github.com/ombima56/insights-edge/internal/util"
)

func main() {
	err := os.Mkdir("db", 0744)
	if err != nil && !os.IsExist(err) {
		log.Printf("Error: %v", err)
	}

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
