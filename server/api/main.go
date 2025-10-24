package api

import (
	"database/sql"
	"fmt"
	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/handlers/auth"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/websockets"
	"log"
	"net/http"
	"time"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// h := &handlers.Handler{DB: db}
	auth := &auth.Handler{DB: db}
	mux.HandleFunc("/ws/stats", websockets.StatWsHandler)
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/auth/signup", auth.SignUpHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/auth/check-setup-status", auth.SetupStatusHandler)
}

func InitApiServer(db *sql.DB) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, db)
	go websockets.BroadcastMetrics()
	handler := middleware.Logger(mux)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
