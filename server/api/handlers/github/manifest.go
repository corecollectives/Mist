package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/handlers"
)

type GithubAppConversion struct {
	ID     int    `json:"id"`
	Slug   string `json:"slug"`
	NodeID string `json:"node_id"`
	Owner  struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"owner"`
	Name          string            `json:"name"`
	ClientID      string            `json:"client_id"`
	ClientSecret  string            `json:"client_secret"`
	WebhookSecret string            `json:"webhook_secret"`
	PEM           string            `json:"pem"`
	HTMLURL       string            `json:"html_url"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ExternalURL   string            `json:"external_url"`
	Permissions   map[string]string `json:"permissions"`
	Events        []string          `json:"events"`
}

func (h *Handler) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing code or state", "code and state are required")
		return
	}

	githubURL := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	req, err := http.NewRequest("POST", githubURL, nil)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create request", err.Error())
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Mist-App")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to make request to GitHub", err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "GitHub API error", string(body))
		return
	}

	var appConversion GithubAppConversion
	if err := json.Unmarshal(body, &appConversion); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to parse GitHub response", err.Error())
		return
	}

	query := `
	INSERT INTO github_app (app_id, client_id, client_secret, webhook_secret, private_key, name, slug, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = h.DB.Exec(
		query,
		appConversion.ID,
		appConversion.ClientID,
		appConversion.ClientSecret,
		appConversion.WebhookSecret,
		appConversion.PEM,
		appConversion.Name,
		appConversion.Slug,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to save GitHub App to database", err.Error())
		return
	}

	redirectUrl := "http://localhost:5173/"
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)

}
