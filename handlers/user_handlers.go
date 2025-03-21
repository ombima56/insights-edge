package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ombima56/insights-edge/database"
	"github.com/ombima56/insights-edge/models"
)

// SignupHandler handles user registration
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var signupReq models.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Received signup request: %+v", signupReq)

	// Check if user exists
	var existingUserID int
	err = database.DB.QueryRow("SELECT id FROM users WHERE wallet_address = ?", signupReq.WalletAddress).Scan(&existingUserID)

	if err == nil {
		// Update existing user
		log.Println("User exists, updating details...")
		_, err = database.DB.Exec(
			`UPDATE users SET email = ?, username = ?, company = ?, business_size = ? WHERE wallet_address = ?`,
			signupReq.Email, signupReq.Username, signupReq.Company, signupReq.BusinessSize, signupReq.WalletAddress,
		)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
			http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
		return
	}

	if err != sql.ErrNoRows {
		log.Printf("Error checking existing user: %v", err)
		http.Error(w, `{"error": "Failed to query user"}`, http.StatusInternalServerError)
		return
	}

	log.Println("New user, inserting into database...")
	result, err := database.DB.Exec(
		`INSERT INTO users (wallet_address, email, username, company, business_size) VALUES (?, ?, ?, ?, ?)`,
		signupReq.WalletAddress, signupReq.Email, signupReq.Username, signupReq.Company, signupReq.BusinessSize,
	)

	if err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, `{"error": "Failed to insert user"}`, http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last inserted ID: %v", err)
	} else {
		log.Printf("User successfully inserted with ID: %d", userID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// GetUsersHandler returns a list of all users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT wallet_address, email, username FROM users")
	if err != nil {
		http.Error(w, `{"error": "Failed to query users"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.WalletAddress, &user.Email, &user.Username)
		if err != nil {
			http.Error(w, `{"error": "Failed to scan user"}`, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
