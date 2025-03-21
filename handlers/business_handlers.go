package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ombima56/insights-edge/contracts"
	"github.com/ombima56/insights-edge/database"
	"github.com/ombima56/insights-edge/models"
)

// RegisterBusinessHandler handles business registration on the blockchain
func RegisterBusinessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from JWT token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse request body
	var registerReq models.BusinessRegistrationRequest
	err = json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if registerReq.Name == "" {
		http.Error(w, `{"error": "Business name is required"}`, http.StatusBadRequest)
		return
	}

	// Get user wallet address
	var walletAddress string
	err = database.DB.QueryRow("SELECT wallet_address FROM users WHERE id = ?", userID).Scan(&walletAddress)
	if err != nil {
		log.Printf("Error getting user wallet address: %v", err)
		http.Error(w, `{"error": "Failed to get user wallet address"}`, http.StatusInternalServerError)
		return
	}

	// Check if user already has a business
	var existingBusinessID int64
	err = database.DB.QueryRow("SELECT id FROM businesses WHERE user_id = ?", userID).Scan(&existingBusinessID)
	if err == nil {
		http.Error(w, `{"error": "User already has a registered business"}`, http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking existing business: %v", err)
		http.Error(w, `{"error": "Failed to check existing business"}`, http.StatusInternalServerError)
		return
	}

	// Insert business record
	result, err := database.DB.Exec(
		"INSERT INTO businesses (user_id, name, industry, location) VALUES (?, ?, ?, ?)",
		userID, registerReq.Name, registerReq.Industry, registerReq.Location,
	)
	if err != nil {
		log.Printf("Error inserting business: %v", err)
		http.Error(w, `{"error": "Failed to register business"}`, http.StatusInternalServerError)
		return
	}

	businessID, _ := result.LastInsertId()

	// In a real implementation, this would trigger a blockchain transaction
	// For now, we'll simulate the blockchain event processing
	contractHandler := contracts.NewSmartContractHandler(database.DB)
	err = contractHandler.ProcessBusinessRegisteredEvent(
		businessID,
		walletAddress,
		registerReq.Name,
		"0x"+generateRandomHex(64), // Simulated transaction hash
	)
	if err != nil {
		log.Printf("Error processing business registration event: %v", err)
		http.Error(w, `{"error": "Failed to process blockchain event"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	response := models.Business{
		ID:       businessID,
		UserID:   userID,
		Name:     registerReq.Name,
		Industry: registerReq.Industry,
		Location: registerReq.Location,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetBusinessProfileHandler returns the business profile for the authenticated user
func GetBusinessProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from JWT token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get business data
	var business models.Business
	err = database.DB.QueryRow(
		`SELECT id, user_id, name, industry, location, created_at 
		FROM businesses WHERE user_id = ?`,
		userID,
	).Scan(
		&business.ID, &business.UserID, &business.Name,
		&business.Industry, &business.Location, &business.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Business not found"}`, http.StatusNotFound)
		} else {
			log.Printf("Error fetching business: %v", err)
			http.Error(w, `{"error": "Failed to fetch business profile"}`, http.StatusInternalServerError)
		}
		return
	}

	// Get subscription data if available
	var subscription models.Subscription
	err = database.DB.QueryRow(
		`SELECT id, business_id, plan_type, start_date, end_date, payment_tx_hash, amount, status 
		FROM subscriptions WHERE business_id = ? AND status = 'active' ORDER BY end_date DESC LIMIT 1`,
		business.ID,
	).Scan(
		&subscription.ID, &subscription.BusinessID, &subscription.PlanType,
		&subscription.StartDate, &subscription.EndDate, &subscription.PaymentTxHash,
		&subscription.Amount, &subscription.Status,
	)

	// Combine business and subscription data
	response := struct {
		Business     models.Business      `json:"business"`
		Subscription *models.Subscription `json:"subscription,omitempty"`
	}{
		Business: business,
	}

	if err == nil {
		response.Subscription = &subscription
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllBusinessesHandler returns a list of all registered businesses
func GetAllBusinessesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	// Optional industry filter
	industry := r.URL.Query().Get("industry")

	var rows *sql.Rows
	var err error

	if industry != "" {
		rows, err = database.DB.Query(
			`SELECT id, user_id, name, industry, location, created_at 
			FROM businesses WHERE industry = ? ORDER BY created_at DESC`,
			industry,
		)
	} else {
		rows, err = database.DB.Query(
			`SELECT id, user_id, name, industry, location, created_at 
			FROM businesses ORDER BY created_at DESC`,
		)
	}

	if err != nil {
		log.Printf("Error querying businesses: %v", err)
		http.Error(w, `{"error": "Failed to fetch businesses"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var businesses []models.Business
	for rows.Next() {
		var business models.Business
		err := rows.Scan(
			&business.ID, &business.UserID, &business.Name,
			&business.Industry, &business.Location, &business.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning business row: %v", err)
			continue
		}
		businesses = append(businesses, business)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(businesses)
}

// PurchaseSubscriptionHandler handles subscription purchases
func PurchaseSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from JWT token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse request body
	var subscriptionReq models.SubscriptionRequest
	err = json.NewDecoder(r.Body).Decode(&subscriptionReq)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if subscriptionReq.PlanType == "" {
		http.Error(w, `{"error": "Plan type is required"}`, http.StatusBadRequest)
		return
	}

	// Get business ID for user
	var businessID int64
	err = database.DB.QueryRow("SELECT id FROM businesses WHERE user_id = ?", userID).Scan(&businessID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "User does not have a registered business"}`, http.StatusBadRequest)
		} else {
			log.Printf("Error getting business ID: %v", err)
			http.Error(w, `{"error": "Failed to get business information"}`, http.StatusInternalServerError)
		}
		return
	}

	// Get user wallet address
	var walletAddress string
	err = database.DB.QueryRow("SELECT wallet_address FROM users WHERE id = ?", userID).Scan(&walletAddress)
	if err != nil {
		log.Printf("Error getting user wallet address: %v", err)
		http.Error(w, `{"error": "Failed to get user wallet address"}`, http.StatusInternalServerError)
		return
	}

	// Calculate subscription price based on plan type
	var amount float64
	var durationDays int64
	switch subscriptionReq.PlanType {
	case "basic":
		amount = 50.0
		durationDays = 30
	case "premium":
		amount = 100.0
		durationDays = 30
	case "enterprise":
		amount = 200.0
		durationDays = 30
	default:
		http.Error(w, `{"error": "Invalid plan type"}`, http.StatusBadRequest)
		return
	}

	// In a real implementation, this would trigger a blockchain transaction
	// For now, we'll simulate the blockchain event processing
	txHash := "0x" + generateRandomHex(64) // Simulated transaction hash
	endDate := durationDays * 24 * 60 * 60 // Duration in seconds

	contractHandler := contracts.NewSmartContractHandler(database.DB)
	err = contractHandler.ProcessSubscriptionPurchasedEvent(
		businessID,
		walletAddress,
		subscriptionReq.PlanType,
		endDate,
		txHash,
		amount,
	)
	if err != nil {
		log.Printf("Error processing subscription purchase event: %v", err)
		http.Error(w, `{"error": "Failed to process blockchain event"}`, http.StatusInternalServerError)
		return
	}

	// Get the created subscription
	var subscription models.Subscription
	err = database.DB.QueryRow(
		`SELECT id, business_id, plan_type, start_date, end_date, payment_tx_hash, amount, status 
		FROM subscriptions WHERE payment_tx_hash = ?`,
		txHash,
	).Scan(
		&subscription.ID, &subscription.BusinessID, &subscription.PlanType,
		&subscription.StartDate, &subscription.EndDate, &subscription.PaymentTxHash,
		&subscription.Amount, &subscription.Status,
	)
	if err != nil {
		log.Printf("Error fetching created subscription: %v", err)
		http.Error(w, `{"error": "Subscription created but failed to fetch details"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}
