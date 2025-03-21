package contracts

import (
	"database/sql"
	"fmt"
)

// ProcessBusinessRegisteredEvent handles a BusinessRegistered event from the smart contract
func (h *SmartContractHandler) ProcessBusinessRegisteredEvent(businessID int64, ownerAddress string, name string, txHash string) error {
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
	err = tx.QueryRow("SELECT id FROM users WHERE wallet_address = ?", ownerAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new user if doesn't exist
			result, err := tx.Exec(
				"INSERT INTO users (wallet_address, username) VALUES (?, ?)",
				ownerAddress,
				fmt.Sprintf("user_%s", ownerAddress[2:10]),
			)
			if err != nil {
				return fmt.Errorf("failed to create user: %v", err)
			}
			userID, err = result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get user ID: %v", err)
			}
		} else {
			return fmt.Errorf("failed to query user: %v", err)
		}
	}
	
	// Insert business
	_, err = tx.Exec(
		"INSERT INTO businesses (id, user_id, name) VALUES (?, ?, ?)",
		businessID,
		userID,
		name,
	)
	if err != nil {
		return fmt.Errorf("failed to insert business: %v", err)
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

// ProcessSubscriptionPurchasedEvent handles a SubscriptionPurchased event
func (h *SmartContractHandler) ProcessSubscriptionPurchasedEvent(businessID int64, subscriber string, planType string, endDate int64, txHash string, amount float64) error {
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
	err = tx.QueryRow("SELECT id FROM users WHERE wallet_address = ?", subscriber).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}
	
	// Record transaction
	_, err = tx.Exec(
		"INSERT INTO transactions (user_id, tx_hash, amount, currency, status) VALUES (?, ?, ?, ?, ?)",
		userID,
		txHash,
		amount,
		"BPT", // Business Platform Token
		"confirmed",
	)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}
	
	// Create subscription record
	_, err = tx.Exec(
		`INSERT INTO subscriptions (business_id, plan_type, start_date, end_date, payment_tx_hash, amount, status) 
         VALUES (?, ?, CURRENT_TIMESTAMP, datetime(?, 'unixepoch'), ?, ?, ?)`,
		businessID,
		planType,
		endDate,
		txHash,
		amount,
		"active",
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %v", err)
	}
	
	// Commit transaction
	return tx.Commit()
}
