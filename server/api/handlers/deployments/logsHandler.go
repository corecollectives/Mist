package deployments

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/websockets"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: websockets.CheckOriginWithSettings,
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
		// Wait for log file to appear
		fileExists := false
		for i := 0; i < 20; i++ {
			if _, err := os.Stat(logPath); err == nil {
				fileExists = true
				break
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
		}

		// If file doesn't exist after waiting, send error event
		if !fileExists {
			select {
			case <-ctx.Done():
				return
			case events <- websockets.DeploymentEvent{
				Type:      "error",
				Timestamp: time.Now(),
				Data: map[string]string{
					"message": "Log file not found. The deployment may have completed without generating logs.",
				},
			}:
			}
			return
		}

		send := make(chan string, 100)
		go func() {
			if err := websockets.WatcherLogs(ctx, logPath, send); err != nil {
				// If error occurs while watching logs, send error event
				select {
				case <-ctx.Done():
				case events <- websockets.DeploymentEvent{
					Type:      "error",
					Timestamp: time.Now(),
					Data: map[string]string{
						"message": "Failed to read log file: " + err.Error(),
					},
				}:
				}
			}
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
