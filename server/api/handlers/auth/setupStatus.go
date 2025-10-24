package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *Handler) SetupStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	db := h.DB
	row := db.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(map[string]bool{"setupRequired": true})
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setupComplete := count > 0
	json.NewEncoder(w).Encode(map[string]bool{"setupRequired": !setupComplete})
}
