package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/corecollectives/mist/fs"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
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
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	err = MigrateDb(db)
	if err != nil {
		return nil, err
	}
	// if err := runMigrations(db); err != nil {
	// 	return nil, fmt.Errorf("failed to run migrations: %v", err)
	// }
	return db, nil
}
