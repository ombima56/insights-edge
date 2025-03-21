package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Database configuration
const (
	dbFile = "business_platform.db"
)

// DB is the global database connection
var DB *sql.DB

// Initialize database and create schema
func InitDB() (*sql.DB, error) {
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

	DB = db
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

	// Read schema.sql file
	schemaPath := filepath.Join("database", "schema.sql")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// Split the schema file into individual statements
	schemaSQL := string(schemaBytes)
	statements := strings.Split(schemaSQL, ";")

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %v\nStatement: %s", err, stmt)
		}
	}

	// Commit transaction
	return tx.Commit()
}
