package models

// User represents a user in the system
type User struct {
	ID            int64  `json:"id,omitempty"`
	WalletAddress string `json:"wallet_address"`
	Email         string `json:"email,omitempty"`
	Username      string `json:"username,omitempty"`
	Company       string `json:"company,omitempty"`
	BusinessSize  string `json:"business_size,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
}

// SignupRequest represents the data needed for user signup
type SignupRequest struct {
	WalletAddress string `json:"wallet_address"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Company       string `json:"company"`
	BusinessSize  string `json:"businessSize"`
}
