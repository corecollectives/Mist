package websockets

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool)
var mu = sync.Mutex{}
var broadcast = make(chan []byte, 5)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			fmt.Println("Client disconnected:", err)
			break
		}
		broadcast <- msg
	}

}

func BroadcastMessages() {
	for {
		message := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func InnitWSServer() {
	http.HandleFunc("/ws", WsHandler)
	go BroadcastMessages()
	fmt.Println("WebSocket server started on :5000")
	http.ListenAndServe(":5000", nil)
}
