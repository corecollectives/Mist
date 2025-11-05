package github

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
)

type RepoListResponse struct {
	TotalCount   int   `json:"total_count"`
	Repositories []any `json:"repositories"`
}

func (h *Handler) GetRepositories(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var installationID string
	err := h.DB.QueryRow(`
		SELECT installation_id FROM github_installations
		WHERE user_id = ?
	`, userData.ID).Scan(&installationID)
	if err == sql.ErrNoRows {
		http.Error(w, "no installation found for user", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("db error: %v", err), http.StatusInternalServerError)
		return
	}

	var (
		token        string
		tokenExpires string
		appID        int
	)
	err = h.DB.QueryRow(`
		SELECT i.access_token, i.token_expires_at, a.app_id
		FROM github_installations i
		JOIN github_app a ON a.id = 1
		WHERE i.installation_id = ?
	`, installationID).Scan(&token, &tokenExpires, &appID)
	if err != nil {
		http.Error(w, "failed to fetch installation info", http.StatusInternalServerError)
		return
	}

	expiry, _ := time.Parse(time.RFC3339, tokenExpires)
	if time.Now().After(expiry) {
		appJWT, err := utils.GenerateGithubJwt(h.DB, appID)
		if err != nil {
			http.Error(w, "failed to generate app jwt", http.StatusInternalServerError)
			return
		}

		newToken, newExpiry, err := regenerateInstallationToken(appJWT, installationID)
		if err != nil {
			http.Error(w, "failed to refresh token", http.StatusInternalServerError)
			return
		}

		_, _ = h.DB.Exec(`
			UPDATE github_installations
			SET access_token = ?, token_expires_at = ?, updated_at = CURRENT_TIMESTAMP
			WHERE installation_id = ?
		`, newToken, newExpiry.Format(time.RFC3339), installationID)

		token = newToken
	}

	allRepos := []any{}
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/installation/repositories?per_page=100&page=%d", page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("GitHub API returned %d", resp.StatusCode), resp.StatusCode)
			return
		}

		var repoList RepoListResponse
		if err := json.NewDecoder(resp.Body).Decode(&repoList); err != nil {
			http.Error(w, "failed to parse GitHub response", http.StatusInternalServerError)
			return
		}

		allRepos = append(allRepos, repoList.Repositories...)

		if len(repoList.Repositories) < 100 {
			break // no more pages
		}
		page++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allRepos)
}

func regenerateInstallationToken(appJWT, installationID string) (string, time.Time, error) {
	url := "https://api.github.com/app/installations/" + installationID + "/access_tokens"

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+appJWT)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", time.Time{}, fmt.Errorf("failed to create token, status %d", resp.StatusCode)
	}

	var tokenResp struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, err
	}

	return tokenResp.Token, tokenResp.ExpiresAt, nil
}
