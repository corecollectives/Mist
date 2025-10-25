package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{"error": message, "success": false, "data": "", "message": message})
}
