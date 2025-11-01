package github

import (
	"database/sql"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetApp(w http.ResponseWriter, r *http.Request) {
	row := h.DB.QueryRow("SELECT id, name, app_id, client_id, slug, created_at FROM github_app LIMIT 1")

	var app models.GitHubApp
	err := row.Scan(&app.ID, &app.Name, &app.AppID, &app.ClientID, &app.Slug, &app.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, true, nil, "GitHub App not found", "no GitHub App configured")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}
	handlers.SendResponse(w, http.StatusOK, true, app, "GitHub App retrieved successfully", "")
}
