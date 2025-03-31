package routes

import (
	"net/http"

	"github.com/ombima56/insights-edge/internal/handlers"
	"github.com/ombima56/insights-edge/internal/middleware"
)

func InitRoutes() {
	// Static file server
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/dashboard", middleware.AuthMiddleware(handlers.DashboardHandler))
	http.HandleFunc("/api/auth/login", handlers.LoginAPIHandler)
	http.HandleFunc("/api/auth/register", handlers.RegisterAPIHandler)
	http.HandleFunc("/api/auth/logout", handlers.LogoutHandler)
}
