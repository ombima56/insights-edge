package models

// Transaction represents a blockchain transaction in the system
type Transaction struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	TxHash    string  `json:"tx_hash"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at,omitempty"`
}
