package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Database configuration
const (
	dbFile = "business_platform.db"
)

type User struct {
	WalletAddress string `json:"wallet_address"`
	Email         string `json:"email"`
	Username      string `json:"username"`
}

var db *sql.DB

// Initialize database and create schema
func initDB() (*sql.DB, error) {
	// Check if database file exists
	_, err := os.Stat(dbFile)
	dbExists := !os.IsNotExist(err)

	// Open database connection
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	// Create tables if database doesn't exist
	if !dbExists {
		log.Println("Creating new database with schema...")
		if err := createSchema(db); err != nil {
			return nil, fmt.Errorf("failed to create schema: %v", err)
		}

		// Insert initial data
		if err := seedInitialData(db); err != nil {
			return nil, fmt.Errorf("failed to seed initial data: %v", err)
		}

		log.Println("Database successfully initialized")
	}

	return db, nil
}

// Create database schema
func createSchema(db *sql.DB) error {
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

	// SQL statements for schema creation
	statements := []string{
		// Users table
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			wallet_address TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE,
			username TEXT UNIQUE,
			company TEXT,               -- Add company column
			business_size TEXT,         -- Add business_size column
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Businesses table
		`CREATE TABLE businesses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			industry TEXT,
			location TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		
		// Market insights table
		`CREATE TABLE market_insights (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			business_id INTEGER,
			industry TEXT NOT NULL,
			insight_type TEXT NOT NULL,
			data TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
		)`,
		
		// Transactions table
		`CREATE TABLE transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			tx_hash TEXT UNIQUE NOT NULL,
			amount REAL NOT NULL,
			currency TEXT NOT NULL,
			status TEXT NOT NULL CHECK(status IN ('pending', 'confirmed', 'failed')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		
		// API access logs table
		`CREATE TABLE api_access_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			endpoint TEXT NOT NULL,
			request_data TEXT,
			response_data TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		
		// Feedback table
		`CREATE TABLE feedback (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			insight_id INTEGER NOT NULL,
			rating INTEGER CHECK(rating BETWEEN 1 AND 5),
			comments TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (insight_id) REFERENCES market_insights(id) ON DELETE CASCADE
		)`,
		
		// Smart contract metadata table
		`CREATE TABLE smart_contracts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			contract_address TEXT UNIQUE NOT NULL,
			contract_name TEXT NOT NULL,
			contract_abi TEXT NOT NULL,
			contract_bytecode TEXT NOT NULL,
			deployed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			owner_user_id INTEGER,
			network TEXT NOT NULL,
			FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE SET NULL
		)`,
		
		// Business subscription payments
		`CREATE TABLE subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			business_id INTEGER NOT NULL,
			plan_type TEXT NOT NULL CHECK(plan_type IN ('basic', 'premium', 'enterprise')),
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			payment_tx_hash TEXT UNIQUE,
			amount REAL NOT NULL,
			status TEXT NOT NULL CHECK(status IN ('active', 'expired', 'canceled')),
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,
			FOREIGN KEY (payment_tx_hash) REFERENCES transactions(tx_hash) ON DELETE SET NULL
		)`,
		
		// Contract events table
		`CREATE TABLE contract_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			contract_id INTEGER NOT NULL,
			event_name TEXT NOT NULL,
			event_data TEXT NOT NULL,
			block_number INTEGER NOT NULL,
			transaction_hash TEXT NOT NULL,
			log_index INTEGER NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (contract_id) REFERENCES smart_contracts(id) ON DELETE CASCADE
		)`,
		
		// Create indices for performance
		`CREATE INDEX idx_users_wallet ON users(wallet_address)`,
		`CREATE INDEX idx_businesses_user ON businesses(user_id)`,
		`CREATE INDEX idx_market_insights_business ON market_insights(business_id)`,
		`CREATE INDEX idx_market_insights_industry ON market_insights(industry)`,
		`CREATE INDEX idx_transactions_user ON transactions(user_id)`,
		`CREATE INDEX idx_transactions_hash ON transactions(tx_hash)`,
		`CREATE INDEX idx_feedback_insight ON feedback(insight_id)`,
		`CREATE INDEX idx_contract_events_contract ON contract_events(contract_id)`,
	}

	// Execute each statement
	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %v\nStatement: %s", err, stmt)
		}
	}

	// Commit transaction
	return tx.Commit()
}

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

// SmartContractHandler handles interactions with blockchain contracts
type SmartContractHandler struct {
	db *sql.DB
}

// Contract represents a smart contract from the database
type Contract struct {
	ID              int    `json:"id"`
	ContractAddress string `json:"contract_address"`
	ContractName    string `json:"contract_name"`
	ContractABI     string `json:"contract_abi"`
	ContractBytecode string `json:"contract_bytecode"`
	DeployedAt      string `json:"deployed_at"`
	OwnerUserID     int    `json:"owner_user_id"`
	Network         string `json:"network"`
}

// GetContractByName retrieves a contract by its name
func (h *SmartContractHandler) GetContractByName(name string) (*Contract, error) {
	var contract Contract
	err := h.db.QueryRow(
		"SELECT id, contract_address, contract_name, contract_abi, contract_bytecode, deployed_at, owner_user_id, network FROM smart_contracts WHERE contract_name = ?",
		name,
	).Scan(
		&contract.ID,
		&contract.ContractAddress,
		&contract.ContractName,
		&contract.ContractABI,
		&contract.ContractBytecode,
		&contract.DeployedAt,
		&contract.OwnerUserID,
		&contract.Network,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %v", err)
	}
	return &contract, nil
}

// LogContractEvent records an event from a smart contract
func (h *SmartContractHandler) LogContractEvent(contractID int, eventName string, eventData map[string]interface{}, blockNumber int64, txHash string, logIndex int) error {
	// Serialize event data to JSON
	eventDataJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}
	
	// Insert event into the database
	_, err = h.db.Exec(
		"INSERT INTO contract_events (contract_id, event_name, event_data, block_number, transaction_hash, log_index) VALUES (?, ?, ?, ?, ?, ?)",
		contractID,
		eventName,
		string(eventDataJSON),
		blockNumber,
		txHash,
		logIndex,
	)
	if err != nil {
		return fmt.Errorf("failed to log contract event: %v", err)
	}
	
	return nil
}

// ProcessBusinessRegisteredEvent handles a BusinessRegistered event from the smart contract
func (h *SmartContractHandler) ProcessBusinessRegisteredEvent(businessID int64, ownerAddress string, name string, txHash string) error {
	// Begin transaction
	tx, err := h.db.Begin()
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
	tx, err := h.db.Begin()
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

// ProcessInsightListedEvent handles an InsightListed event
func (h *SmartContractHandler) ProcessInsightListedEvent(insightID int64, provider string, industry string, insightType string, dataHash string, price float64, txHash string) error {
	// Begin transaction
	tx, err := h.db.Begin()
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
	tx, err := h.db.Begin()
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

func signupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
        return
    }

    var user struct {
        WalletAddress string `json:"wallet_address"`
        Email         string `json:"email"`
        Username      string `json:"username"`
        Company       string `json:"company"`
        BusinessSize  string `json:"businessSize"`
    }

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    log.Printf("Received signup request: %+v", user)

    // Check if user exists
    var existingUserID int
    err = db.QueryRow("SELECT id FROM users WHERE wallet_address = ?", user.WalletAddress).Scan(&existingUserID)

    if err == nil {
        // Update existing user
        log.Println("User exists, updating details...")
        _, err = db.Exec(
            `UPDATE users SET email = ?, username = ?, company = ?, business_size = ? WHERE wallet_address = ?`,
            user.Email, user.Username, user.Company, user.BusinessSize, user.WalletAddress,
        )
        if err != nil {
            log.Printf("Failed to update user: %v", err)
            http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
        return
    }

    if err != sql.ErrNoRows {
        log.Printf("Error checking existing user: %v", err)
        http.Error(w, `{"error": "Failed to query user"}`, http.StatusInternalServerError)
        return
    }

    log.Println("New user, inserting into database...")
    result, err := db.Exec(
        `INSERT INTO users (wallet_address, email, username, company, business_size) VALUES (?, ?, ?, ?, ?)`,
        user.WalletAddress, user.Email, user.Username, user.Company, user.BusinessSize,
    )

    if err != nil {
        log.Printf("Error inserting user: %v", err)
        http.Error(w, `{"error": "Failed to insert user"}`, http.StatusInternalServerError)
        return
    }

    userID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error retrieving last inserted ID: %v", err)
    } else {
        log.Printf("User successfully inserted with ID: %d", userID)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT wallet_address, email, username FROM users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query users: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.WalletAddress, &user.Email, &user.Username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan user: %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Serve the index.html file from the templates directory
	http.ServeFile(w, r, "frontend/templates/index.html")
}

func main() {
	// Initialize database
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve the home page
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/users", getUsersHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}