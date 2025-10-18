package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	return db, nil
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	rows, err := db.Query("SELECT version from schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan applied migration version: %w", err)
		}
		appliedMigrations[version] = true
	}

	migrationsDir := "db/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory '%s': %w", migrationsDir, err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) != ".sql" {
			continue
		}
		if appliedMigrations[filename] {
			continue
		}
		log.Println("applying migration: %s\n", filename)

		filepath := filepath.Join(migrationsDir, filename)
		content, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read migration file '%s': %w", filepath, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration '%s': %w", filename, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration '%s': %w", filename, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", filename); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration '%s': %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for migration '%s': %w", filename, err)
		}

	}
	log.Println("All migrations applied successfully.")
	return nil
}
