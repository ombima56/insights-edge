package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	_ = os.Mkdir("db", 0755)

	db, err := sql.Open("sqlite3", "./db/insights.db")
	if err != nil {
		return err
	}

	DB = db

	// Create users table if it doesn't exist
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			wallet_addr TEXT UNIQUE NOT NULL,
			first_name TEXT,
			last_name TEXT,
			company_name TEXT,
			industry TEXT,
			account_type TEXT NOT NULL DEFAULT 'free',
			created_at TEXT NOT NULL
		)
	`); err != nil {
		return err
	}

	// Create insights table if it doesn't exist
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS insights (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider TEXT NOT NULL,
			industry TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			price TEXT NOT NULL,
			created_at TEXT NOT NULL
		)
	`); err != nil {
		return err
	}

	// Create purchases table if it doesn't exist
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS purchases (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			insight_id INTEGER NOT NULL,
			buyer TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY (insight_id) REFERENCES insights(id)
		)
	`); err != nil {
		return err
	}

	return nil
}
