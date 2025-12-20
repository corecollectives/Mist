package deployments

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/websockets"
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

		log.Printf("Deployment logs WebSocket rejected from origin: %s (allowed: %s)", origin, settings.AllowedOrigins)
		return false
	},
}

func LogsHandler(w http.ResponseWriter, r *http.Request) {
	depIdstr := r.URL.Query().Get("id")
	depId, err := strconv.ParseInt(depIdstr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "invalid deployment id", err.Error())
		return
	}
	dep, err := models.GetDeploymentByID(depId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "deployment not found", err.Error())
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to upgrade to websocket", err.Error())
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	logPath := docker.GetLogsPath(dep.CommitHash, depId)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	events := make(chan websockets.DeploymentEvent, 100)

	go websockets.WatchDeploymentStatus(ctx, depId, events)

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					cancel()
					return
				}
			}
		}
	}()

	go func() {
		for i := 0; i < 20; i++ {
			if _, err := os.Stat(logPath); err == nil {
				break
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
		}

		send := make(chan string, 100)
		go func() {
			_ = websockets.WatcherLogs(ctx, logPath, send)
			close(send)
		}()

		for line := range send {
			select {
			case <-ctx.Done():
				return
			case events <- websockets.DeploymentEvent{
				Type:      "log",
				Timestamp: time.Now(),
				Data: websockets.LogUpdate{
					Line:      line,
					Timestamp: time.Now(),
				},
			}:
			}
		}
	}()

	for event := range events {
		msg, err := json.Marshal(event)
		if err != nil {
			continue
		}

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			cancel()
			break
		}

		if statusData, ok := event.Data.(websockets.StatusUpdate); ok {
			if statusData.Status == "success" || statusData.Status == "failed" {
				time.Sleep(1 * time.Second)
				cancel()
				break
			}
		}
	}
}
