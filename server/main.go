package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/middleware"
)

func main() {
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)
	handler := middleware.CORSMiddleware(middleware.Logger(mux))
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
