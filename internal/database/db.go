package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./insights.db")
	if err != nil {
		return err
	}

	// Create users table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			account_type TEXT NOT NULL,
			company_name TEXT,
			industry TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create sessions table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create market_insights table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS market_insights (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			industry TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			trend_value REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Insert some dummy data for market insights
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO market_insights (industry, title, description, trend_value)
		VALUES 
		('Technology', 'Cloud Computing Growth', 'Enterprise cloud adoption increased by 25% in 2024', 25.5),
		('Technology', 'AI Implementation', 'AI integration in business processes grew by 40%', 40.0),
		('Technology', 'Cybersecurity Trends', 'Investment in cybersecurity solutions up by 30%', 30.0),
		('Retail', 'E-commerce Expansion', 'Online retail sales increased by 35% globally', 35.0),
		('Retail', 'Mobile Shopping', 'Mobile commerce accounts for 60% of online sales', 60.0),
		('Retail', 'Sustainable Packaging', 'Green packaging adoption up by 45%', 45.0),
		('Healthcare', 'Telemedicine Growth', 'Virtual healthcare consultations up by 50%', 50.0),
		('Healthcare', 'AI Diagnostics', 'AI-powered diagnostic accuracy improved by 28%', 28.0),
		('Healthcare', 'Remote Monitoring', 'Remote patient monitoring adoption up by 42%', 42.0)
	`)

	return err
}
