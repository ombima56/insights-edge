package models

import "errors"

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // Don't expose password in JSON
	WalletAddr   string    `json:"wallet_addr"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	AccountType  string    `json:"account_type"`
	CompanyName  string    `json:"company_name"`
	Industry     string    `json:"industry"`
	CreatedAt    string    `json:"created_at"`
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents the request body for registration
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	WalletAddr string `json:"wallet_addr"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AccountType string `json:"account_type"`
	CompanyName string `json:"company_name"`
	Industry   string `json:"industry"`
}

// Error definitions
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidToken      = errors.New("invalid token")
)

// Insight represents a market insight
type Insight struct {
	ID          int64     `json:"id"`
	Provider    string    `json:"provider"`
	Industry    string    `json:"industry"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	CreatedAt   string    `json:"created_at"`
}

// Purchase represents a purchase of an insight
type Purchase struct {
	ID        int64     `json:"id"`
	InsightID int64     `json:"insight_id"`
	Buyer     string    `json:"buyer"`
	CreatedAt string    `json:"created_at"`
}
