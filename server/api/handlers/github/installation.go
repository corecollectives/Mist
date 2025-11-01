package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/utils"
)

type InstallationTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type InstallationInfo struct {
	ID      int64 `json:"id"`
	AppID   int64 `json:"app_id"`
	Account struct {
		Login string `json:"login"`
		Type  string `json:"type"`
		ID    int64  `json:"id"`
	} `json:"account"`
}

func (h *Handler) HandleInstallationEvent(w http.ResponseWriter, r *http.Request) {
	baseFrontendURL := GetFrontendBaseUrl()
	installationId := r.URL.Query().Get("installation_id")
	state := r.URL.Query().Get("state")

	if installationId == "" || state == "" {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=missing_params", baseFrontendURL), http.StatusSeeOther)
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

	appJWT, err := utils.GenerateGithubJwt(h.DB, stateData.AppId)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=failed_to_generate_jwt", baseFrontendURL), http.StatusSeeOther)
		return
	}

	tokenReq, _ := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationId), nil)
	tokenReq.Header.Set("Authorization", "Bearer "+appJWT)
	tokenReq.Header.Set("Accept", "application/vnd.github+json")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil || tokenResp.StatusCode != http.StatusCreated {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=failed_to_create_installation_token", baseFrontendURL), http.StatusSeeOther)
		return
	}
	defer tokenResp.Body.Close()

	var token InstallationTokenResponse
	if err := json.NewDecoder(tokenResp.Body).Decode(&token); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=failed_to_parse_token", baseFrontendURL), http.StatusSeeOther)
		return
	}

	infoReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/app/installations/%s", installationId), nil)
	infoReq.Header.Set("Authorization", "Bearer "+appJWT)
	infoReq.Header.Set("Accept", "application/vnd.github+json")

	infoResp, err := http.DefaultClient.Do(infoReq)
	if err != nil || infoResp.StatusCode != http.StatusOK {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=failed_to_fetch_installation_info", baseFrontendURL), http.StatusSeeOther)
		return
	}
	defer infoResp.Body.Close()

	var installInfo InstallationInfo
	if err := json.NewDecoder(infoResp.Body).Decode(&installInfo); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=failed_to_parse_installation_info", baseFrontendURL), http.StatusSeeOther)
		return
	}

	_, err = h.DB.Exec(`
		INSERT OR REPLACE INTO github_installations
		(installation_id, account_login, account_type, access_token, token_expires_at, user_id, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, installInfo.ID, installInfo.Account.Login, installInfo.Account.Type, token.Token, token.ExpiresAt, stateData.UserId)

	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=failed_to_store_installation_info", baseFrontendURL), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/callback?toast=Github_App_Created_Successfully&redirect=/git", baseFrontendURL), http.StatusSeeOther)
}
