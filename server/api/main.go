package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/handlers/auth"
	"github.com/corecollectives/mist/api/handlers/projects"
	"github.com/corecollectives/mist/api/handlers/users"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/websockets"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	h := &handlers.Handler{DB: db}
	auth := &auth.Handler{DB: db}
	proj := &projects.Handler{DB: db}
	users := &users.Handler{DB: db}
	mux.Handle("/api/ws/stats", middleware.AuthMiddleware(h)(http.HandlerFunc(websockets.StatWsHandler)))
	mux.HandleFunc("/api/health", handlers.HealthCheckHandler)

	mux.HandleFunc("/api/auth/signup", auth.SignUpHandler)
	mux.HandleFunc("/api/auth/login", auth.LoginHandler)
	mux.HandleFunc("/api/auth/me", auth.MeHandler)
	mux.HandleFunc("/api/auth/logout", auth.LogoutHandler)

	mux.Handle("/api/users/create", middleware.AuthMiddleware(h)(http.HandlerFunc(users.CreateUser)))

	mux.Handle("/api/projects/create", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.CreateProject)))
	mux.Handle("/api/projects/getAll", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.GetProjects)))
	mux.Handle("/api/projects/getFromId", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.GetProjectFromId)))
	mux.Handle("/api/projects/update", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.UpdateProject)))
	mux.Handle("/api/projects/delete", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.DeleteProject)))

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
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
