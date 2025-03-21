package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/ombima56/insights-edge/database"
	"golang.org/x/crypto/bcrypt"
)

// Hash a password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare a password with a hash
func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Parse JSON from request body
func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// Respond with JSON
func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Get database connection
func getDB() *sql.DB {
	return database.DB
}

// Check if user is a business
func isBusinessUser(userID int64) (bool, error) {
	var isBusiness bool
	err := getDB().QueryRow("SELECT is_business FROM users WHERE id = ?", userID).Scan(&isBusiness)
	return isBusiness, err
}

// Get business ID for a user
func getBusinessIDForUser(userID int64) (int64, error) {
	var businessID int64
	err := getDB().QueryRow("SELECT id FROM businesses WHERE user_id = ?", userID).Scan(&businessID)
	return businessID, err
}

// Check if user has purchased an insight
func hasUserPurchasedInsight(userID, insightID int64) (bool, error) {
	var count int
	err := getDB().QueryRow(
		"SELECT COUNT(*) FROM insight_purchases WHERE user_id = ? AND insight_id = ?",
		userID, insightID,
	).Scan(&count)
	return count > 0, err
}

// Check if user has an active subscription
func hasActiveSubscription(businessID int64) (bool, error) {
	var count int
	err := getDB().QueryRow(
		"SELECT COUNT(*) FROM subscriptions WHERE business_id = ? AND status = 'active' AND end_date > CURRENT_TIMESTAMP",
		businessID,
	).Scan(&count)
	return count > 0, err
}

// Get user wallet address
func getUserWalletAddress(userID int64) (string, error) {
	var walletAddress string
	err := getDB().QueryRow("SELECT wallet_address FROM users WHERE id = ?", userID).Scan(&walletAddress)
	return walletAddress, err
}

// Format error response
func formatErrorResponse(message string) string {
	return `{"error": "` + message + `"}`
}

// Generate a random hex string (if not using the one from user_handlers.go)
func generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	return bytes, err
}

// Convert bytes to hex string
func bytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
