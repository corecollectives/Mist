package main

import (
	"fmt"
	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/api/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)
	handler := middleware.Logger(mux)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
