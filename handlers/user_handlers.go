package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"
	"strings"

	"github.com/ombima56/insights-edge/database"
	"github.com/ombima56/insights-edge/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret for development - in production, this should be set via environment variable
		return "insights-edge-secret-key-2025"
	}
	return secret
}

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

	// Validate required fields
	if signupReq.Email == "" || signupReq.Username == "" || signupReq.Password == "" || signupReq.UserType == "" {
		http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		return
	}

	// Validate user type
	if signupReq.UserType != "individual" && signupReq.UserType != "business" {
		http.Error(w, `{"error": "Invalid user type. Must be 'individual' or 'business'"}`, http.StatusBadRequest)
		return
	}

	// Check if user exists by email
	var existingUserID int
	err = database.DB.QueryRow("SELECT id FROM users WHERE email = ?", signupReq.Email).Scan(&existingUserID)
	if err == nil {
		http.Error(w, `{"error": "Email already in use"}`, http.StatusConflict)
		return
	}

	// Generate password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, `{"error": "Failed to process password"}`, http.StatusInternalServerError)
		return
	}

	// Generate a wallet address if not provided
	walletAddress := signupReq.WalletAddress
	if walletAddress == "" {
		walletAddress = "0x" + generateRandomHex(40) // Simple placeholder for now
	}

	// Insert new user
	result, err := database.DB.Exec(
		`INSERT INTO users (wallet_address, email, username, company, business_size, user_type, password_hash) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		walletAddress, signupReq.Email, signupReq.Username, signupReq.Company, 
		signupReq.BusinessSize, signupReq.UserType, string(hashedPassword),
	)

	if err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()

	// Create JWT token
	token, err := generateJWT(userID, signupReq.Email, signupReq.UserType)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, `{"error": "Failed to generate authentication token"}`, http.StatusInternalServerError)
		return
	}

	// Return user data and token
	user := models.User{
		ID:            userID,
		WalletAddress: walletAddress,
		Email:         signupReq.Email,
		Username:      signupReq.Username,
		Company:       signupReq.Company,
		BusinessSize:  signupReq.BusinessSize,
		UserType:      signupReq.UserType,
	}

	response := models.AuthResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LoginHandler handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var loginReq models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Find user by email
	var user models.User
	var passwordHash string
	
	err = database.DB.QueryRow(
		`SELECT id, wallet_address, email, username, company, business_size, user_type, password_hash 
		FROM users WHERE email = ?`, 
		loginReq.Email,
	).Scan(
		&user.ID, &user.WalletAddress, &user.Email, &user.Username, 
		&user.Company, &user.BusinessSize, &user.UserType, &passwordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		} else {
			log.Printf("Database error during login: %v", err)
			http.Error(w, `{"error": "Authentication failed"}`, http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Email, user.UserType)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, `{"error": "Failed to generate authentication token"}`, http.StatusInternalServerError)
		return
	}

	// Return user data and token
	response := models.AuthResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserProfileHandler returns the authenticated user's profile
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Fetch user data
	var user models.User
	err = database.DB.QueryRow(
		`SELECT id, wallet_address, email, username, company, business_size, user_type, created_at 
		FROM users WHERE id = ?`, 
		userID,
	).Scan(
		&user.ID, &user.WalletAddress, &user.Email, &user.Username, 
		&user.Company, &user.BusinessSize, &user.UserType, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		} else {
			log.Printf("Database error fetching user profile: %v", err)
			http.Error(w, `{"error": "Failed to fetch user profile"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUsersHandler returns a list of all users (admin only)
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, wallet_address, email, username, company, business_size, user_type, created_at 
		FROM users`,
	)
	if err != nil {
		http.Error(w, `{"error": "Failed to query users"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.WalletAddress, &user.Email, &user.Username, 
			&user.Company, &user.BusinessSize, &user.UserType, &user.CreatedAt,
		)
		if err != nil {
			http.Error(w, `{"error": "Failed to scan user"}`, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Helper functions

// Claims represents the JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}

// generateJWT creates a new JWT token for a user
func generateJWT(userID int64, email, userType string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// getUserIDFromToken extracts the user ID from the JWT token
func getUserIDFromToken(r *http.Request) (int64, error) {
	tokenString := extractTokenFromHeader(r)
	if tokenString == "" {
		return 0, jwt.ErrSignatureInvalid
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	return claims.UserID, nil
}

// extractTokenFromHeader gets the JWT token from the Authorization header
func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}
	return bearerToken
}

// generateRandomHex generates a random hex string of the specified length
func generateRandomHex(length int) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = hexChars[time.Now().UnixNano()%16]
		time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
	}
	return string(result)
}
