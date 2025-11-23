package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type EnvVariable struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"appId"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func CreateEnvVariable(appID int64, key, value string) (*EnvVariable, error) {
	id := utils.GenerateRandomId()
	query := `
		INSERT INTO envs (id, app_id, key, value, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, app_id, key, value, created_at, updated_at
	`
	var env EnvVariable
	err := db.QueryRow(query, id, appID, key, value).Scan(
		&env.ID, &env.AppID, &env.Key, &env.Value, &env.CreatedAt, &env.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func GetEnvVariablesByAppID(appID int64) ([]EnvVariable, error) {
	query := `
		SELECT id, app_id, key, value, created_at, updated_at
		FROM envs
		WHERE app_id = ?
		ORDER BY key ASC
	`
	rows, err := db.Query(query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []EnvVariable
	for rows.Next() {
		var env EnvVariable
		err := rows.Scan(&env.ID, &env.AppID, &env.Key, &env.Value, &env.CreatedAt, &env.UpdatedAt)
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	return envs, nil
}

func UpdateEnvVariable(id int64, key, value string) error {
	query := `
		UPDATE envs
		SET key = ?, value = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, key, value, id)
	return err
}

func DeleteEnvVariable(id int64) error {
	query := `DELETE FROM envs WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func GetEnvVariableByID(id int64) (*EnvVariable, error) {
	query := `
		SELECT id, app_id, key, value, created_at, updated_at
		FROM envs
		WHERE id = ?
	`
	var env EnvVariable
	err := db.QueryRow(query, id).Scan(&env.ID, &env.AppID, &env.Key, &env.Value, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &env, nil
}
