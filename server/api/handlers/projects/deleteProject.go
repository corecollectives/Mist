package projects

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
)

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	projectIdStr := r.URL.Query().Get("id")
	if projectIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing project ID", "project id is required")
		return
	}
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid project ID", err.Error())
		return
	}

	var ownerId int64
	err = h.DB.QueryRow("SELECT owner_id FROM projects WHERE id = ?", projectId).Scan(&ownerId)
	if err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "no such project")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}
	if ownerId != userData.ID {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized to delete this project", "permission denied")
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to begin transaction", err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM project_members WHERE project_id = ?", projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete project members", err.Error())
		return
	}

	res, err := tx.Exec("DELETE FROM projects WHERE id = ?", projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete project", err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Error getting affected rows", err.Error())
		return
	}
	if rowsAffected == 0 {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "no such project")
		return
	}

	if err := tx.Commit(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to commit transaction", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, nil, "Project deleted successfully", "")
}
