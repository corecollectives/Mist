package github

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"os"
)

type Handler struct {
	DB *sql.DB
}

type StateData struct {
	AppId  int `json:"appId"`
	UserId int `json:"userId"`
}

func GetFrontendBaseUrl() string {
	if os.Getenv("ENV") == "dev" {
		return "http://localhost:5173"
	}
	return ""
}

func GenerateState(appId int, userId int) string {
	payload := map[string]int{
		"appId":  appId,
		"userId": userId,
	}
	jsonBytes, _ := json.Marshal(payload)
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return encoded
}

func CheckIfAppExists(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM github_app`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
