package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"insights-edge/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	AccountType string `json:"accountType"`
	CompanyName string `json:"companyName,omitempty"`
	Industry    string `json:"industry,omitempty"`
	Password    string `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Password    string `json:"password"`
	AccountType string `json:"accountType"`
	CompanyName string `json:"companyName,omitempty"`
	Industry    string `json:"industry,omitempty"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login", nil)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "register", nil)
		return
	}
}

func LoginAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", req.Email).Scan(&user.ID, &user.Email, &hashedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	token := generateSessionToken()
	expiresAt := time.Now().Add(24 * time.Hour)
	_, err = database.DB.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", user.ID, token, expiresAt)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func RegisterAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (email, first_name, last_name, password_hash, account_type, company_name, industry) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	result, err := database.DB.Exec(query, 
		req.Email, 
		req.FirstName, 
		req.LastName, 
		hashedPassword, 
		req.AccountType, 
		req.CompanyName, 
		req.Industry)
	
	if err != nil {
		http.Error(w, "Email already exists or invalid data", http.StatusBadRequest)
		return
	}

	userID, _ := result.LastInsertId()

	// Create session
	token := generateSessionToken()
	expiresAt := time.Now().Add(24 * time.Hour)
	_, err = database.DB.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusCreated)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil {
		database.DB.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func generateSessionToken() string {
	// In a production environment, use a proper random token generator
	return "session_" + time.Now().Format("20060102150405")
}
