package deployments

import (
	"context"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"

	"github.com/corecollectives/mist/websockets"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
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

	logPath := docker.GetLogsPath(dep.CommitHash, depId)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	send := make(chan string)
	go func() {
		_ = websockets.WatcherLogs(ctx, logPath, send)
		close(send)
	}()

	for line := range send {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			cancel()
			break
		}
	}

}
