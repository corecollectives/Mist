package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/websockets"
)

func RegisterRoutes(mux *http.ServeMux) {
	//jo protected routes honge usme middleware use kar lena jaise /dashboard waghera jo bhi
	mux.HandleFunc("/ws", websockets.WsHandler)
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/signup", handlers.SignUpHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/doesExist", handlers.DoesUserExist)
}

func InitApiServer() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	go websockets.BroadcastMessages() //need to run this goroutine before starting the server to handle broadcasting.
	handler := middleware.Logger(mux)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
