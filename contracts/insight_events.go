package contracts

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// ProcessInsightListedEvent handles an InsightListed event
func (h *SmartContractHandler) ProcessInsightListedEvent(insightID int64, provider string, industry string, insightType string, dataHash string, price float64, txHash string) error {
	// Begin transaction
	tx, err := h.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	
	// Get user by wallet address
	var userID int64
	err = tx.QueryRow("SELECT id FROM users WHERE wallet_address = ?", provider).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}
	
	// Get business ID for user (use first business if multiple)
	var businessID sql.NullInt64
	err = tx.QueryRow("SELECT id FROM businesses WHERE user_id = ? LIMIT 1", userID).Scan(&businessID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query business: %v", err)
	}
	
	// Create insight data structure
	insightData := map[string]interface{}{
		"dataHash": dataHash,
		"price":    price,
		"provider": provider,
	}
	
	insightDataJSON, err := json.Marshal(insightData)
	if err != nil {
		return fmt.Errorf("failed to marshal insight data: %v", err)
	}
	
	// Insert market insight
	_, err = tx.Exec(
		"INSERT INTO market_insights (id, business_id, industry, insight_type, data) VALUES (?, ?, ?, ?, ?)",
		insightID,
		sql.NullInt64{Int64: businessID.Int64, Valid: businessID.Valid},
		industry,
		insightType,
		string(insightDataJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert market insight: %v", err)
	}
	
	// Record transaction
	_, err = tx.Exec(
		"INSERT INTO transactions (user_id, tx_hash, amount, currency, status) VALUES (?, ?, ?, ?, ?)",
		userID,
		txHash,
		0, // Gas cost not tracked
		"ETH",
		"confirmed",
	)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}
	
	// Commit transaction
	return tx.Commit()
}

// ProcessInsightPurchasedEvent handles an InsightPurchased event
func (h *SmartContractHandler) ProcessInsightPurchasedEvent(insightID int64, buyer string, seller string, price float64, txHash string) error {
	// Begin transaction
	tx, err := h.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	
	// Get buyer user by wallet address
	var buyerUserID int64
	err = tx.QueryRow("SELECT id FROM users WHERE wallet_address = ?", buyer).Scan(&buyerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new user if doesn't exist
			result, err := tx.Exec(
				"INSERT INTO users (wallet_address, username) VALUES (?, ?)",
				buyer,
				fmt.Sprintf("user_%s", buyer[2:10]),
			)
			if err != nil {
				return fmt.Errorf("failed to create user: %v", err)
			}
			buyerUserID, err = result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get user ID: %v", err)
			}
		} else {
			return fmt.Errorf("failed to query user: %v", err)
		}
	}
	
	// Record transaction
	_, err = tx.Exec(
		"INSERT INTO transactions (user_id, tx_hash, amount, currency, status) VALUES (?, ?, ?, ?, ?)",
		buyerUserID,
		txHash,
		price,
		"BPT", // Business Platform Token
		"confirmed",
	)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}
	
	// Commit transaction
	return tx.Commit()
}
