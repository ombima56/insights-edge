package models

// MarketInsight represents a market insight in the system
type MarketInsight struct {
	ID         int64  `json:"id"`
	BusinessID int64  `json:"business_id,omitempty"`
	Industry   string `json:"industry"`
	InsightType string `json:"insight_type"`
	Data       string `json:"data"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// Feedback represents user feedback on a market insight
type Feedback struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	InsightID int64  `json:"insight_id"`
	Rating    int    `json:"rating"`
	Comments  string `json:"comments,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
