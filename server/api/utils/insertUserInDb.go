package utils

import (
	"database/sql"

	"github.com/corecollectives/mist/models"
	"golang.org/x/crypto/bcrypt"
)

func InsertUserInDb(db *sql.DB, username, email, password, role string) (models.User, error) {
	var user models.User
	id := GenerateRandomId()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	tx, err := db.Begin()
	if err != nil {
		return user, err
	}

	err = tx.QueryRow(`
		INSERT INTO users (id, username, email, password_hash, role)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, username, email, role, created_at, updated_at
	`, id, username, email, passwordHash, role).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return user, err
	}

	if err := tx.Commit(); err != nil {
		return user, err
	}

	return user, nil
}
