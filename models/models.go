package models

// BusinessRegistrationRequest represents a request to register a business
type BusinessRegistrationRequest struct {
	Name     string `json:"name"`
	Industry string `json:"industry"`
	Location string `json:"location"`
}

// SubscriptionRequest represents a request to purchase a subscription
type SubscriptionRequest struct {
	PlanType string  `json:"plan_type"`
	Amount   float64 `json:"amount"`
}

// InsightRequest represents a request to list a new insight
type InsightRequest struct {
	Industry    string  `json:"industry"`
	InsightType string  `json:"insight_type"`
	Title       string  `json:"title"`
	Data        string  `json:"data"`
	Price       float64 `json:"price"`
}

// FeedbackRequest represents a request to submit feedback
type FeedbackRequest struct {
	Rating   int    `json:"rating"`
	Comments string `json:"comments"`
}
