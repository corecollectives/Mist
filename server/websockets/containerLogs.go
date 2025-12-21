package websockets

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ContainerLogsEvent struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

var containerLogsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOriginWithSettings,
}

func ContainerLogsHandler(w http.ResponseWriter, r *http.Request) {
	appIDStr := r.URL.Query().Get("appId")
	if appIDStr == "" {
		http.Error(w, "appId is required", http.StatusBadRequest)
		return
	}

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid appId", http.StatusBadRequest)
		return
	}

	app, err := models.GetApplicationByID(appID)
	if err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	conn, err := containerLogsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade websocket connection for container logs")
		return
	}
	defer conn.Close()

	log.Info().Int64("app_id", appID).Str("app_name", app.Name).Msg("Container logs client connected")

	containerName := docker.GetContainerName(app.Name, appID)

	if !docker.ContainerExists(containerName) {
		conn.WriteJSON(ContainerLogsEvent{
			Type:      "error",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"message": "Container not found",
			},
		})
		return
	}

	status, err := docker.GetContainerStatus(containerName)
	if err != nil {
		conn.WriteJSON(ContainerLogsEvent{
			Type:      "error",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"message": fmt.Sprintf("Failed to get container status: %v", err),
			},
		})
		return
	}

	conn.WriteJSON(ContainerLogsEvent{
		Type:      "status",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"container": containerName,
			"state":     status.State,
			"status":    status.Status,
		},
	})

	if status.State != "running" {
		conn.WriteJSON(ContainerLogsEvent{
			Type:      "error",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"message": fmt.Sprintf("Container is not running (state: %s)", status.State),
			},
		})
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logChan := make(chan string, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(logChan)

		cmd := exec.CommandContext(ctx, "sh", "-c",
			fmt.Sprintf("docker logs -f --tail 100 %s 2>&1", containerName))

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errChan <- fmt.Errorf("failed to create stdout pipe: %w", err)
			return
		}

		if err := cmd.Start(); err != nil {
			errChan <- fmt.Errorf("failed to start docker logs: %w", err)
			return
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
				return
			case logChan <- line:
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("scanner error: %w", err)
		}

		cmd.Wait()
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
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Info().Int64("app_id", appID).Msg("Container logs client disconnected")
				}
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int64("app_id", appID).Msg("Container logs context cancelled")
			return

		case err := <-errChan:
			log.Error().Err(err).Int64("app_id", appID).Msg("Container logs stream error")
			conn.WriteJSON(ContainerLogsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": err.Error(),
				},
			})
			return

		case line, ok := <-logChan:
			if !ok {
				conn.WriteJSON(ContainerLogsEvent{
					Type:      "end",
					Timestamp: time.Now().Format(time.RFC3339),
					Data: map[string]interface{}{
						"message": "Log stream ended",
					},
				})
				return
			}

			event := ContainerLogsEvent{
				Type:      "log",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"line": line,
				},
			}

			if err := conn.WriteJSON(event); err != nil {
				log.Warn().Err(err).Int64("app_id", appID).Msg("Failed to send container log message to client")
				return
			}
		}
	}
}
