package api

import (
	"github.com/corecollectives/mist/api/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
}
