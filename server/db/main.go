package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	dbPath := ""
	if os.Getenv("ENV") == "dev" {
		dbPath = "./mist.db"
	} else {
		dbPath = "/var/lib/mist/mist.db"
	}
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %v", err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	fmt.Println("Database path:", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	//just for dev to see where file is created as db init is lazy.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS dummy (id INTEGER PRIMARY KEY)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create dummy table: %v", err)
	}
	return db, nil
}
