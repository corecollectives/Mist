package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GithubAppConversion struct {
	ID            int               `json:"id"`
	Slug          string            `json:"slug"`
	NodeID        string            `json:"node_id"`
	Name          string            `json:"name"`
	ClientID      string            `json:"client_id"`
	ClientSecret  string            `json:"client_secret"`
	WebhookSecret string            `json:"webhook_secret"`
	PEM           string            `json:"pem"`
	HTMLURL       string            `json:"html_url"`
	ExternalURL   string            `json:"external_url"`
	Permissions   map[string]string `json:"permissions"`
	Events        []string          `json:"events"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Owner         struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"owner"`
}

func (h *Handler) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	baseFrontendURL := GetFrontendBaseUrl()

	if code == "" || state == "" {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=missing_code_or_state", baseFrontendURL), http.StatusSeeOther)
		return
	}

	decodedState, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=invalid_state_encoding", baseFrontendURL), http.StatusSeeOther)
		return
	}

	var stateData StateData
	if err := json.Unmarshal(decodedState, &stateData); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=invalid_state_data", baseFrontendURL), http.StatusSeeOther)
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code), nil)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=request_creation_failed", baseFrontendURL), http.StatusSeeOther)
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Mist-App")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=github_request_failed", baseFrontendURL), http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=github_api_error&details=%s", baseFrontendURL, base64.URLEncoding.EncodeToString(body)), http.StatusSeeOther)
		return
	}

	var app GithubAppConversion
	if err := json.Unmarshal(body, &app); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=invalid_github_response", baseFrontendURL), http.StatusSeeOther)
		return
	}

	_, err = h.DB.Exec(`
		INSERT INTO github_app 
			(app_id, client_id, client_secret, webhook_secret, private_key, name, slug, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, app.ID, app.ClientID, app.ClientSecret, app.WebhookSecret, app.PEM, app.Name, app.Slug, time.Now(), time.Now())
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=db_insert_failed", baseFrontendURL), http.StatusSeeOther)
		return
	}

	newState := GenerateState(app.ID, stateData.UserId)
	redirectURL := fmt.Sprintf("https://github.com/apps/%s/installations/new?state=%s", app.Slug, newState)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
