package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckOriginWithSettings,
}

var StatClients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func StatWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	mu.Lock()
	StatClients[conn] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(StatClients, conn)
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
		if len(StatClients) == 0 {
			continue
		}
		metrics, err := GetStats()
		if err != nil {
			log.Printf("Error getting metrics: %v", err)
			continue
		}
		msg, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("Error marshalling metrics: %v", err)
			continue
		}
		mu.Lock()
		for client := range StatClients {
			client.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Close()
				delete(StatClients, client)
			}
		}
		mu.Unlock()
	}
}
