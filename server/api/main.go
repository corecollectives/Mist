package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/handlers/auth"
	"github.com/corecollectives/mist/api/handlers/github"
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
	github := &github.Handler{DB: db}
	mux.Handle("/api/ws/stats", middleware.AuthMiddleware(h)(http.HandlerFunc(websockets.StatWsHandler)))
	mux.HandleFunc("GET /api/health", handlers.HealthCheckHandler)

	mux.HandleFunc("POST /api/auth/signup", auth.SignUpHandler)
	mux.HandleFunc("POST /api/auth/login", auth.LoginHandler)
	mux.HandleFunc("GET /api/auth/me", auth.MeHandler)
	mux.HandleFunc("POST /api/auth/logout", auth.LogoutHandler)

	mux.Handle("POST /api/users/create", middleware.AuthMiddleware(h)(http.HandlerFunc(users.CreateUser)))
	mux.Handle("GET /api/users/getAll", middleware.AuthMiddleware(h)(http.HandlerFunc(users.GetUsers)))
	mux.Handle("GET /api/users/getFromId", middleware.AuthMiddleware(h)(http.HandlerFunc(users.GetUserById)))
	mux.Handle("DELETE /api/users/delete", middleware.AuthMiddleware(h)(http.HandlerFunc(users.DeleteUser)))

	mux.Handle("POST /api/projects/create", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.CreateProject)))
	mux.Handle("GET /api/projects/getAll", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.GetProjects)))
	mux.Handle("GET /api/projects/getFromId", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.GetProjectFromId)))
	mux.Handle("PUT /api/projects/update", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.UpdateProject)))
	mux.Handle("DELETE /api/projects/delete", middleware.AuthMiddleware(h)(http.HandlerFunc(proj.DeleteProject)))

	mux.Handle("GET /api/github/app", middleware.AuthMiddleware(h)(http.HandlerFunc(github.GetApp)))
	mux.Handle("GET /api/github/app/create", middleware.AuthMiddleware(h)(http.HandlerFunc(github.CreateGithubApp)))
	mux.Handle("GET /api/github/callback", http.HandlerFunc(github.CallBackHandler))
	mux.Handle("GET /api/github/installation/callback", http.HandlerFunc(github.HandleInstallationEvent))
	mux.Handle("GET /api/github/repositories", middleware.AuthMiddleware(h)(http.HandlerFunc(github.GetRepositories)))

}

func InitApiServer(db *sql.DB) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, db)
	staticDir := "static"
	fs := http.FileServer(http.Dir(staticDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, r.URL.Path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})
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
