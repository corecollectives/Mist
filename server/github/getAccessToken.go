package github

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"
)

// Structs for DB and GitHub API
type GithubApp struct {
	AppID      int
	PrivateKey string
}

type GithubInstallation struct {
	InstallationID int
	AccessToken    sql.NullString
	TokenExpiresAt sql.NullTime
}

// Function: GetGitHubAccessToken
func GetGitHubAccessToken(db *sql.DB, userID int) (string, error) {
	var app GithubApp
	var inst GithubInstallation

	// 1️⃣ Fetch app credentials
	err := db.QueryRow(`SELECT app_id, private_key FROM github_app LIMIT 1`).Scan(
		&app.AppID,
		&app.PrivateKey,
	)
	if err != nil {
		return "", fmt.Errorf("failed to fetch github app credentials: %w", err)
	}

	// 2️⃣ Fetch installation details for this user
	err = db.QueryRow(`
		SELECT installation_id, access_token, token_expires_at
		FROM github_installations WHERE user_id = ?
	`, userID).Scan(&inst.InstallationID, &inst.AccessToken, &inst.TokenExpiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to fetch github installation: %w", err)
	}

	// 3️⃣ Check if existing token is valid
	if inst.AccessToken.Valid && inst.TokenExpiresAt.Valid && time.Until(inst.TokenExpiresAt.Time) > 5*time.Minute {
		return inst.AccessToken.String, nil
	}

	// 4️⃣ Create JWT for GitHub App
	jwt, err := createGitHubJWT(app.AppID, app.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create JWT: %w", err)
	}

	// 5️⃣ Request installation access token
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", inst.InstallationID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request installation token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		var body bytes.Buffer
		body.ReadFrom(resp.Body)
		return "", fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, body.String())
	}

	var result struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	// 6️⃣ Update DB with new token
	_, _ = db.Exec(`
		UPDATE github_installations
		SET access_token = ?, token_expires_at = ?
		WHERE user_id = ?
	`, result.Token, result.ExpiresAt, userID)

	return result.Token, nil
}

// Helper: Create JWT for GitHub App
func createGitHubJWT(appID int, privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("invalid private key PEM")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse key: %w", err)
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	now := time.Now().Unix()
	payload := fmt.Sprintf(`{"iat":%d,"exp":%d,"iss":%d}`, now-60, now+540, appID)
	payloadEnc := base64.RawURLEncoding.EncodeToString([]byte(payload))

	toSign := header + "." + payloadEnc
	hash := sha256.Sum256([]byte(toSign))

	sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	sigEnc := base64.RawURLEncoding.EncodeToString(sig)
	return toSign + "." + sigEnc, nil
}
