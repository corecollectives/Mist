package github

import (
	"database/sql"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetApp(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok || userInfo == nil {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized", "user not authenticated")
		return
	}

	query := `
		SELECT 
			a.id,
			a.name,
			a.app_id,
			a.client_id,
			a.slug,
			a.created_at,
			CASE WHEN i.installation_id IS NOT NULL THEN 1 ELSE 0 END AS is_installed
		FROM github_app a
		LEFT JOIN github_installations i ON i.user_id = ?
		WHERE a.id = 1
	`

	row := h.DB.QueryRow(query, userInfo.ID)

	var app models.GitHubApp
	var isInstalled bool

	err := row.Scan(
		&app.ID,
		&app.Name,
		&app.AppID,
		&app.ClientID,
		&app.Slug,
		&app.CreatedAt,
		&isInstalled,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, true, nil, "GitHub App not found", "")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"app":         app,
		"isInstalled": isInstalled,
	}, "GitHub App retrieved successfully", "")
}
