package api

import (
	"database/sql"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/handlers/applications"
	"github.com/corecollectives/mist/api/handlers/auth"
	"github.com/corecollectives/mist/api/handlers/deployments"
	"github.com/corecollectives/mist/api/handlers/github"
	"github.com/corecollectives/mist/api/handlers/projects"
	"github.com/corecollectives/mist/api/handlers/users"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/websockets"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	proj := &projects.Handler{DB: db}

	mux.Handle("/api/ws/stats", middleware.AuthMiddleware()(http.HandlerFunc(websockets.StatWsHandler)))
	mux.HandleFunc("GET /api/health", handlers.HealthCheckHandler)

	mux.HandleFunc("POST /api/auth/signup", auth.SignUpHandler)
	mux.HandleFunc("POST /api/auth/login", auth.LoginHandler)
	mux.HandleFunc("GET /api/auth/me", auth.MeHandler)
	mux.HandleFunc("POST /api/auth/logout", auth.LogoutHandler)

	mux.Handle("POST /api/users/create", middleware.AuthMiddleware()(http.HandlerFunc(users.CreateUser)))
	mux.Handle("GET /api/users/getAll", middleware.AuthMiddleware()(http.HandlerFunc(users.GetUsers)))
	mux.Handle("GET /api/users/getFromId", middleware.AuthMiddleware()(http.HandlerFunc(users.GetUserById)))
	mux.Handle("DELETE /api/users/delete", middleware.AuthMiddleware()(http.HandlerFunc(users.DeleteUser)))

	mux.Handle("POST /api/projects/create", middleware.AuthMiddleware()(http.HandlerFunc(proj.CreateProject)))
	mux.Handle("GET /api/projects/getAll", middleware.AuthMiddleware()(http.HandlerFunc(proj.GetProjects)))
	mux.Handle("GET /api/projects/getFromId", middleware.AuthMiddleware()(http.HandlerFunc(proj.GetProjectFromId)))
	mux.Handle("PUT /api/projects/update", middleware.AuthMiddleware()(http.HandlerFunc(proj.UpdateProject)))
	mux.Handle("DELETE /api/projects/delete", middleware.AuthMiddleware()(http.HandlerFunc(proj.DeleteProject)))
	mux.Handle("POST /api/projects/addMember", middleware.AuthMiddleware()(http.HandlerFunc(proj.AddMember)))

	mux.Handle("POST /api/apps/create", middleware.AuthMiddleware()(http.HandlerFunc(applications.CreateApplication)))
	mux.Handle("POST /api/apps/getByProjectId", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetApplicationByProjectID)))
	mux.Handle("POST /api/apps/getById", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetApplicationById)))
	mux.Handle("PUT /api/apps/update", middleware.AuthMiddleware()(http.HandlerFunc(applications.UpdateApplication)))
	mux.Handle("POST /api/apps/getLatestCommit", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetLatestCommit)))

	mux.Handle("GET /api/github/app", middleware.AuthMiddleware()(http.HandlerFunc(github.GetApp)))
	mux.Handle("GET /api/github/app/create", middleware.AuthMiddleware()(http.HandlerFunc(github.CreateGithubApp)))
	mux.Handle("GET /api/github/callback", http.HandlerFunc(github.CallBackHandler))
	mux.Handle("GET /api/github/installation/callback", http.HandlerFunc(github.HandleInstallationEvent))
	mux.Handle("GET /api/github/repositories", middleware.AuthMiddleware()(http.HandlerFunc(github.GetRepositories)))
	mux.Handle("POST /api/github/branches", middleware.AuthMiddleware()(http.HandlerFunc(github.GetBranches)))
	mux.HandleFunc("POST /api/github/webhook", github.GithubWebhook)

	mux.HandleFunc("/api/ws/logs", deployments.LogsHandler)
	mux.Handle("POST /api/deployments/create", middleware.AuthMiddleware()(http.HandlerFunc(deployments.AddDeployHandler)))
	mux.Handle("POST /api/deployments/getByAppId", middleware.AuthMiddleware()(http.HandlerFunc(deployments.GetByApplicationID)))

}
