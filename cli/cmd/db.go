package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/corecollectives/mist/models"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "/var/lib/mist/mist.db"

// initDB initializes the database connection for CLI
// It expects the database file to already exist
func initDB() error {
	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found at %s. Please ensure Mist is installed and running", dbPath)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set the database instance for models
	models.SetDB(db)
	return nil
}
