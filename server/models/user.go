package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	AvatarURL    *string   `json:"avatarUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (u *User) Create() error {
	query := `
		INSERT INTO users (id, username, email, password_hash, role, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, role, avatar_url, created_at, updated_at
	`
	u.ID = utils.GenerateRandomId()
	err := db.QueryRow(query, u.ID, u.Username, u.Email, u.PasswordHash, u.Role, u.AvatarURL).Scan(
		&u.ID, &u.Username, &u.Email, &u.Role, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	return err
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func GetUserByID(userID int64) (*User, error) {
	query := `
	  SELECT id, username, email, role, avatar_url, created_at, updated_at
	  FROM users
	  WHERE id = $1
	`
	user := &User{}
	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUserByID(userID int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, userID)
	return err
}

func UpdateUser(u *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2,  role = $3, avatar_url = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`
	return db.QueryRow(query, u.Username, u.Email, u.Role, u.AvatarURL, u.ID).Scan(&u.UpdatedAt)
}

func (u *User) MatchPassword(password string) bool {
	query := `
		SELECT password_hash
		FROM users
		WHERE id = $1
	`
	var storedHash string
	err := db.QueryRow(query, u.ID).Scan(&storedHash)
	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) == nil
}

func (u *User) UpdatePassword() error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
  `
	return db.QueryRow(query, u.PasswordHash, u.ID).Scan(&u.UpdatedAt)

}

func GetUserByEmail(email string) (*User, error) {
	query := `
	  SELECT id, username, email, role, avatar_url, created_at, updated_at
	  FROM users
	  WHERE email = $1
	`
	user := &User{}
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	query := `
	  SELECT id, username, email, role, avatar_url, created_at, updated_at
	  FROM users
	  WHERE username = $1
	`
	user := &User{}
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserPassword(userID int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := db.Exec(query, passwordHash, userID)
	return err
}

func GetAllUsers() ([]User, error) {
	query := `
		SELECT id, username, email, role, avatar_url, created_at, updated_at
		FROM users
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserRole(userID int64) (string, error) {
	query := `
		SELECT role
		FROM users
		WHERE id = $1
	`
	var role string
	err := db.QueryRow(query, userID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func GetUserCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
