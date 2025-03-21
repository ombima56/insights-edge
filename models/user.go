package models

// User represents a user in the system
type User struct {
	ID            int64  `json:"id,omitempty"`
	WalletAddress string `json:"wallet_address"`
	Email         string `json:"email,omitempty"`
	Username      string `json:"username,omitempty"`
	Company       string `json:"company,omitempty"`
	BusinessSize  string `json:"business_size,omitempty"`
	UserType      string `json:"user_type,omitempty"`
	PasswordHash  string `json:"-"`
	CreatedAt     string `json:"created_at,omitempty"`
}

// SignupRequest represents the data needed for user signup
type SignupRequest struct {
	WalletAddress string `json:"wallet_address"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Company       string `json:"company"`
	BusinessSize  string `json:"businessSize"`
	UserType      string `json:"userType"`
	Password      string `json:"password"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token,omitempty"`
}
