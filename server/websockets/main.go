package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		fmt.Println("Upgrade error:", err)
		return
	}
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

}

func BroadcastMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if len(clients) == 0 {
			continue
		}
		metrics, err := GetMetrics()
		if err != nil {
			log.Println("Error in getting the metrics: ", err)
			continue
		}
		msg, err := json.Marshal(metrics)
		if err != nil {
			log.Println("error in marshallung metrics: ", err)
			continue
		}
		mu.Lock()
		for client := range clients {
			client.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("Error sending message, removing client:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}
