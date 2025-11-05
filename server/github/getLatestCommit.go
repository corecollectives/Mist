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

type LatestCommit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	URL     string `json:"html_url"`
	Author  string `json:"author"`
}

func GetLatestCommit(db *sql.DB, appID, userID int64) (*LatestCommit, error) {
	var repoName, branch string
	err := db.QueryRow(`SELECT git_repository, COALESCE(git_branch, 'main') FROM apps WHERE id = ?`, appID).
		Scan(&repoName, &branch)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app repo: %w", err)
	}
	if repoName == "" {
		return nil, fmt.Errorf("app has no linked GitHub repository")
	}

	var installationID int64
	err = db.QueryRow(`SELECT installation_id FROM github_installations WHERE user_id = ?`, userID).
		Scan(&installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch installation: %w", err)
	}

	var appAppID int64
	var privateKey string
	err = db.QueryRow(`SELECT app_id, private_key FROM github_app LIMIT 1`).Scan(&appAppID, &privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub App credentials: %w", err)
	}

	jwtToken, err := generateAppJWTManual(appAppID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub App JWT: %w", err)
	}

	accessToken, err := getInstallationAccessToken(installationID, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %w", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s", repoName, branch)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var data struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commit"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub response: %w", err)
	}

	return &LatestCommit{
		SHA:     data.SHA,
		Message: data.Commit.Message,
		URL:     data.HTMLURL,
		Author:  data.Commit.Author.Name,
	}, nil
}

func generateAppJWTManual(appID int64, privateKeyPEM string) (string, error) {
	now := time.Now().Unix()
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	}
	payload := map[string]interface{}{
		"iat": now - 60,
		"exp": now + 600,
		"iss": appID,
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	encode := func(b []byte) string {
		return base64.RawURLEncoding.EncodeToString(b)
	}

	headerEncoded := encode(headerJSON)
	payloadEncoded := encode(payloadJSON)

	data := headerEncoded + "." + payloadEncoded

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("invalid PEM private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse RSA key: %w", err)
	}

	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	signed := data + "." + encode(signature)
	return signed, nil
}

func getInstallationAccessToken(installationID int64, jwtToken string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return "", fmt.Errorf("GitHub token request failed (%s): %s", resp.Status, buf.String())
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("invalid token response: %w", err)
	}

	return result.Token, nil
}
