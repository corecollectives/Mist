package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		// Get system settings from database
		settings, err := models.GetSystemSettings()
		if err != nil {
			log.Printf("ERROR: Failed to get system settings for WebSocket CORS: %v", err)
			// Deny by default if we can't get settings
			return false
		}

		// In non-production mode, allow localhost
		if !settings.ProductionMode {
			if strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				origin == "" {
				return true
			}
		}

		// Check against allowed origins (comma-separated list)
		if settings.AllowedOrigins != "" {
			allowedList := strings.Split(settings.AllowedOrigins, ",")
			for _, allowed := range allowedList {
				allowed = strings.TrimSpace(allowed)
				if origin == allowed {
					return true
				}
			}
		}

		// Allow empty origin (same-origin requests)
		if origin == "" {
			return true
		}

		log.Printf("WebSocket connection rejected from origin: %s (allowed: %s)", origin, settings.AllowedOrigins)
		return false
	},
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
