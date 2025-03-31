package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ombima56/insights-edge/backend/database"
	"github.com/ombima56/insights-edge/backend/handlers"
	"github.com/ombima56/insights-edge/backend/repository"
	"github.com/ombima56/insights-edge/backend/routes"
	"github.com/ombima56/insights-edge/backend/service"
)

func main() {
	dbErr := database.InitDB()
	if dbErr != nil {
		log.Fatal("Failed to initialize database:", dbErr)
	}
	defer func() {
		if err := database.DB.Close(); err != nil {
			log.Println("Failed to close database connection:", err)
		}
	}()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	userRepository := repository.NewUserRepository(database.DB)
	insightRepository := repository.NewInsightRepository(database.DB)

	userService := service.NewUserService(userRepository)
	insightService := service.NewInsightService(insightRepository, userRepository)

	authHandler := handlers.NewAuthHandler(userService, store)
	insightHandler := handlers.NewInsightHandler(insightService)

	r := mux.NewRouter()

	routes.SetupRoutes(r, authHandler, insightHandler, store)

	srv := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}

	log.Printf("Server started on http://localhost%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
