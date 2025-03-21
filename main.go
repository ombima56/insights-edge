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

	// Register HTTP handlers for pages
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	// http.HandleFunc("/marketplace", handlers.MarketplaceHandler)
	// http.HandleFunc("/insight/", handlers.InsightDetailHandler)

	// User API endpoints
	http.HandleFunc("/api/signup", handlers.SignupHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/user/profile", handlers.GetUserProfileHandler)
	// http.HandleFunc("/api/user/password", handlers.UpdatePasswordHandler)
	http.HandleFunc("/api/users", handlers.GetUsersHandler)

	// Business API endpoints
	http.HandleFunc("/api/business/register", handlers.RegisterBusinessHandler)
	http.HandleFunc("/api/business/profile", handlers.GetBusinessProfileHandler)
	http.HandleFunc("/api/business/all", handlers.GetAllBusinessesHandler)
	// http.HandleFunc("/api/business/subscription", handlers.GetBusinessSubscriptionHandler)
	http.HandleFunc("/api/business/subscription/purchase", handlers.PurchaseSubscriptionHandler)

	// Insight API endpoints
	// http.HandleFunc("/api/insights", handlers.GetInsightsHandler)
	// http.HandleFunc("/api/insights/list", handlers.ListInsightHandler)
	// http.HandleFunc("/api/insights/purchased", handlers.GetPurchasedInsightsHandler)
	// http.HandleFunc("/api/insights/", handlers.GetInsightByIDHandler)
	// http.HandleFunc("/api/insights/purchase", handlers.PurchaseInsightHandler)
	// http.HandleFunc("/api/insights/feedback", handlers.SubmitFeedbackHandler)

	// Transaction API endpoints
	// http.HandleFunc("/api/transactions", handlers.GetTransactionsHandler)
	// http.HandleFunc("/api/dashboard/stats", handlers.GetDashboardStatsHandler)

	// Start the server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
