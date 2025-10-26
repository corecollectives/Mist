package handlers

import (
	"database/sql"
	"errors"

	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetUserFromId(userID int64) (*models.User, error) {
	var user models.User

	err := h.DB.QueryRow(`
        SELECT id, username, email, password_hash, role, created_at, updated_at
        FROM users WHERE id = ?`, userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
