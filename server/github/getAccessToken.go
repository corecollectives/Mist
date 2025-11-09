package github

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/corecollectives/mist/models"
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
func GetGitHubAccessToken(userID int) (string, error) {

	app, _, err := models.GetApp(userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch github app credentials: %w", err)
	}

	// 2️⃣ Fetch installation details for this user
	inst, err := models.GetInstallationByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch github installation: %w", err)
	}

	// 3️⃣ Check if existing token is valid
	if time.Until(inst.TokenExpiresAt) > 5*time.Minute {
		return inst.AccessToken, nil
	}

	// 4️⃣ Create JWT for GitHub App
	jwt, err := GenerateGithubJwt(int(app.AppID))
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
	err = models.UpdateInstallationToken(inst.InstallationID, result.Token, result.ExpiresAt)

	return result.Token, nil
}
