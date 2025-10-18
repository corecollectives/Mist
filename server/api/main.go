package api

import (
	"fmt"
	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"log"
	"net/http"
	"time"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
}

func InitApiServer() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := middleware.Logger(mux)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
