package api

import (
	"database/sql"
	"fmt"
	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/websockets"
	"log"
	"net/http"
	"time"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	h := &handlers.Handler{DB: db}
	//jo protected routes honge usme middleware use kar lena jaise /dashboard waghera jo bhi
	mux.HandleFunc("/ws/stats", websockets.StatWsHandler)
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/signup", h.SignUpHandler)
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/doesExist", h.DoesUserExist)
}

func InitApiServer(db *sql.DB) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, db)
	go websockets.BroadcastMetrics() //need to run this goroutine before starting the server to handle broadcasting.
	handler := middleware.Logger(mux)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
