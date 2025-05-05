package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ombima56/insights-edge/internal/database"
)   

// GetRelevantInsightsHandler returns insights relevant to the user's business interest
func GetRelevantInsightsHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// For individual users, get insights based on their business interest
	if user.AccountType == "individual" && user.BusinessInterest != "" {
		query := `
            SELECT id, industry, title, description, type, content, user_id, created_at 
            FROM market_insights 
            WHERE industry = ? 
            ORDER BY created_at DESC
        `
		rows, err := database.DB.Query(query, user.BusinessInterest)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var insights []MarketInsight
		for rows.Next() {
			var insight MarketInsight
			err := rows.Scan(&insight.ID, &insight.Industry, &insight.Title, &insight.Description, 
				&insight.Type, &insight.Content, &insight.UserID, &insight.CreatedAt)
			if err != nil {
				continue
			}
			insights = append(insights, insight)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(insights)
		return
	}

	// If no business interest is set or user is not individual
	http.Error(w, "No business interest specified", http.StatusBadRequest)
}

// GetBusinessInsightsHandler returns insights created by the business user
func GetBusinessInsightsHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only business users can access this endpoint
	if user.AccountType != "business" {
		http.Error(w, "Only business users can access this endpoint", http.StatusForbidden)
		return
	}

	query := `
		SELECT id, industry, title, description, type, content, user_id, created_at 
		FROM market_insights 
		WHERE user_id = ? 
		ORDER BY created_at DESC
	`
	rows, err := database.DB.Query(query, user.ID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var insights []MarketInsight
	for rows.Next() {
		var insight MarketInsight
		err := rows.Scan(&insight.ID, &insight.Industry, &insight.Title, &insight.Description, 
			&insight.Type, &insight.Content, &insight.UserID, &insight.CreatedAt)
		if err != nil {
			continue
		}
		insights = append(insights, insight)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insights)
}

// CreateInsightHandler creates a new insight
func CreateInsightHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only business users can create insights
	if user.AccountType != "business" {
		http.Error(w, "Only business users can create insights", http.StatusForbidden)
		return
	}

	var insightData struct {
		Title       string `json:"title"`
		Industry    string `json:"industry"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&insightData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if insightData.Title == "" || insightData.Industry == "" || 
		insightData.Type == "" || insightData.Content == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Insert the insight
	query := `
		INSERT INTO market_insights (title, industry, type, description, content, user_id) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := database.DB.Exec(query, 
		insightData.Title, 
		insightData.Industry,
		insightData.Type, 
		insightData.Description, 
		insightData.Content, 
		user.ID)
	
	if err != nil {
		http.Error(w, "Failed to create insight", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly created insight
	insightID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get insight ID", http.StatusInternalServerError)
		return
	}

	// Return the created insight
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      insightID,
		"message": "Insight created successfully",
	})
}

// UpdateInsightHandler updates an existing insight
func UpdateInsightHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only business users can update insights
	if user.AccountType != "business" {
		http.Error(w, "Only business users can update insights", http.StatusForbidden)
		return
	}

	// Get insight ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/insights/")
	insightID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid insight ID", http.StatusBadRequest)
		return
	}

	// Check if the insight belongs to the user
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM market_insights WHERE id = ? AND user_id = ?", 
		insightID, user.ID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "Insight not found or not owned by you", http.StatusNotFound)
		return
	}

	var updateData struct {
		Title       string `json:"title"`
		Industry    string `json:"industry"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if updateData.Title == "" || updateData.Industry == "" || 
		updateData.Type == "" || updateData.Content == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Update the insight
	query := `
		UPDATE market_insights 
		SET title = ?, industry = ?, type = ?, description = ?, content = ? 
		WHERE id = ? AND user_id = ?
	`
	_, err = database.DB.Exec(query, 
		updateData.Title, 
		updateData.Industry,
		updateData.Type, 
		updateData.Description, 
		updateData.Content, 
		insightID, 
		user.ID)
	
	if err != nil {
		http.Error(w, "Failed to update insight", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Insight updated successfully",
	})
}
