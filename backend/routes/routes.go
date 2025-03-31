package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ombima56/insights-edge/backend/handlers"
	"github.com/ombima56/insights-edge/backend/middleware"
)

func SetupRoutes(r *mux.Router, authHandler *handlers.AuthHandler, insightHandler *handlers.InsightHandler, store *sessions.CookieStore) {
	// Static file handler for assets
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	// Public routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET")
	r.HandleFunc("/api/auth/register", authHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.LoginUser).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddlewareWithStore(store))

	protected.HandleFunc("/insights", insightHandler.GetInsights).Methods("GET")
	protected.HandleFunc("/insights", insightHandler.CreateInsight).Methods("POST")
	protected.HandleFunc("/insights/{id}", insightHandler.GetInsight).Methods("GET")
	protected.HandleFunc("/insights/purchase", insightHandler.PurchaseInsight).Methods("POST")
	protected.HandleFunc("/insights/my-insights", insightHandler.GetMyInsights).Methods("GET")
	protected.HandleFunc("/insights/my-purchases", insightHandler.GetMyPurchases).Methods("GET")

	// Public dashboard route (authentication handled by the handler)
	r.HandleFunc("/dashboard", handlers.DashboardHandler).Methods("GET")
}
