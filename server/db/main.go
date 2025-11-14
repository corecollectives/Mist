package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/corecollectives/mist/fs"
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
	err := fs.CreateDirIfNotExists(dbDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	return db, nil
}
