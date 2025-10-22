package handlers

import (
	"log"
	"net/http"
)

func (h *Handler) SetupStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SetupStatusHandler called")
	db := h.DB
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	userCount, err := db.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userCount.Next() {
		var count int
		err = userCount.Scan(&count)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			w.Write([]byte(`{"setup_complete": true}`))
		}
	} else {
		w.Write([]byte(`{"setup_complete": false}`))
	}

}
