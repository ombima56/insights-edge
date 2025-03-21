package models

// Business represents a business entity in the system
type Business struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	Industry  string `json:"industry,omitempty"`
	Location  string `json:"location,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Subscription represents a business subscription
type Subscription struct {
	ID            int64  `json:"id"`
	BusinessID    int64  `json:"business_id"`
	PlanType      string `json:"plan_type"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	PaymentTxHash string `json:"payment_tx_hash,omitempty"`
	Amount        float64 `json:"amount"`
	Status        string `json:"status"`
}
