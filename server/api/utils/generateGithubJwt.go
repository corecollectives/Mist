package utils

import (
	"database/sql"
	"encoding/pem"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateGithubJwt(db *sql.DB, appID int) (string, error) {
	var appNumericId string
	var appPrivateKeyPEM string

	err := db.QueryRow(`
		SELECT app_id, private_key FROM github_app WHERE id = 1
	`).Scan(&appNumericId, &appPrivateKeyPEM)

	if err != nil {
		return "", err
	}

	block, _ := pem.Decode([]byte(appPrivateKeyPEM))
	if block == nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(appPrivateKeyPEM))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": appNumericId,
	})

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
