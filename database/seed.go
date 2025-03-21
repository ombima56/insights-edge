package database

import (
	"database/sql"
	"fmt"
)

// Seed initial data
func seedInitialData(db *sql.DB) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Insert admin user
	adminResult, err := tx.Exec(
		`INSERT INTO users (wallet_address, email, username) VALUES (?, ?, ?)`,
		"0x1234567890AbCdEf1234567890AbCdEf12345678",
		"admin@businessplatform.io",
		"admin",
	)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %v", err)
	}
	
	adminID, err := adminResult.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get admin user ID: %v", err)
	}
	
	// Insert platform smart contracts
	businessRegistryABI := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"businessId","type":"uint256"},{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"string","name":"name","type":"string"}],"name":"BusinessRegistered","type":"event"},{"inputs":[{"internalType":"string","name":"_name","type":"string"},{"internalType":"string","name":"_industry","type":"string"},{"internalType":"string","name":"_location","type":"string"}],"name":"registerBusiness","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"}]`
	
	subscriptionManagerABI := `[{"inputs":[{"internalType":"address","name":"_token","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"businessId","type":"uint256"},{"indexed":true,"internalType":"address","name":"subscriber","type":"address"},{"indexed":false,"internalType":"string","name":"planType","type":"string"},{"indexed":false,"internalType":"uint256","name":"endDate","type":"uint256"}],"name":"SubscriptionPurchased","type":"event"},{"inputs":[{"internalType":"uint256","name":"_businessId","type":"uint256"},{"internalType":"string","name":"_planType","type":"string"},{"internalType":"uint256","name":"_duration","type":"uint256"}],"name":"purchaseSubscription","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	
	insightMarketplaceABI := `[{"inputs":[{"internalType":"address","name":"_token","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"insightId","type":"uint256"},{"indexed":true,"internalType":"address","name":"provider","type":"address"},{"indexed":false,"internalType":"string","name":"industry","type":"string"},{"indexed":false,"internalType":"uint256","name":"price","type":"uint256"}],"name":"InsightListed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"insightId","type":"uint256"},{"indexed":true,"internalType":"address","name":"buyer","type":"address"},{"indexed":true,"internalType":"address","name":"seller","type":"address"},{"indexed":false,"internalType":"uint256","name":"price","type":"uint256"}],"name":"InsightPurchased","type":"event"},{"inputs":[{"internalType":"string","name":"_industry","type":"string"},{"internalType":"string","name":"_insightType","type":"string"},{"internalType":"string","name":"_dataHash","type":"string"},{"internalType":"uint256","name":"_price","type":"uint256"}],"name":"listInsight","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_insightId","type":"uint256"}],"name":"purchaseInsight","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	
	tokenContractABI := `[{"inputs":[{"internalType":"string","name":"_name","type":"string"},{"internalType":"string","name":"_symbol","type":"string"},{"internalType":"uint256","name":"_initialSupply","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
	
	// Insert token contract
	_, err = tx.Exec(
		`INSERT INTO smart_contracts (contract_address, contract_name, contract_abi, contract_bytecode, owner_user_id, network) VALUES (?, ?, ?, ?, ?, ?)`,
		"0xA1B2C3D4E5F6A7B8C9D0E1F2A3B4C5D6E7F8A9B0",
		"BusinessPlatformToken",
		tokenContractABI,
		"0x60806040...", // placeholder for actual bytecode
		adminID,
		"ethereum",
	)
	if err != nil {
		return fmt.Errorf("failed to insert token contract: %v", err)
	}
	
	// Insert business registry contract
	_, err = tx.Exec(
		`INSERT INTO smart_contracts (contract_address, contract_name, contract_abi, contract_bytecode, owner_user_id, network) VALUES (?, ?, ?, ?, ?, ?)`,
		"0xB1C2D3E4F5G6H7I8J9K0L1M2N3O4P5Q6R7S8T9",
		"BusinessRegistry",
		businessRegistryABI,
		"0x60806040...", // placeholder for actual bytecode
		adminID,
		"ethereum",
	)
	if err != nil {
		return fmt.Errorf("failed to insert business registry contract: %v", err)
	}
	
	// Insert subscription manager contract
	_, err = tx.Exec(
		`INSERT INTO smart_contracts (contract_address, contract_name, contract_abi, contract_bytecode, owner_user_id, network) VALUES (?, ?, ?, ?, ?, ?)`,
		"0xC1D2E3F4G5H6I7J8K9L0M1N2O3P4Q5R6S7T8U9",
		"SubscriptionManager",
		subscriptionManagerABI,
		"0x60806040...", // placeholder for actual bytecode
		adminID,
		"ethereum",
	)
	if err != nil {
		return fmt.Errorf("failed to insert subscription manager contract: %v", err)
	}
	
	// Insert insight marketplace contract
	_, err = tx.Exec(
		`INSERT INTO smart_contracts (contract_address, contract_name, contract_abi, contract_bytecode, owner_user_id, network) VALUES (?, ?, ?, ?, ?, ?)`,
		"0xD1E2F3G4H5I6J7K8L9M0N1O2P3Q4R5S6T7U8V9",
		"InsightMarketplace",
		insightMarketplaceABI,
		"0x60806040...", // placeholder for actual bytecode
		adminID,
		"ethereum",
	)
	if err != nil {
		return fmt.Errorf("failed to insert insight marketplace contract: %v", err)
	}

	// Commit transaction
	return tx.Commit()
}
