package dockerdeploy

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"

	// "github.com/corecollectives/mist/queue"
	"github.com/corecollectives/mist/websockets"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Deployer struct {
	DB           *sql.DB
	LogDirectory string
}

// var deployer *Deployer

func (d *Deployer) LogsHandler(w http.ResponseWriter, r *http.Request) {
	depIdstr := r.URL.Query().Get("id")
	depId, err := strconv.ParseInt(depIdstr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "invalid deployment id", err.Error())
		return
	}
	dep, err := d.loadDeployment(depId)
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

	logPath := d.GetLogsPath(dep.CommitHash)
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

func (d *Deployer) loadDeployment(depId int64) (*models.Deployment, error) {
	row := d.DB.QueryRow("SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at FROM deployments WHERE id = ?", depId)
	dep := &models.Deployment{}
	var triggeredBy sql.NullInt64
	var finishedAt sql.NullTime
	var logs sql.NullString
	err := row.Scan(&dep.ID, &dep.AppID, &dep.CommitHash, &dep.CommitMessage, &triggeredBy, &logs, &dep.Status, &dep.CreatedAt, &finishedAt)
	if err != nil {
		return nil, err
	}
	if triggeredBy.Valid {
		dep.TriggeredBy = &triggeredBy.Int64
	}
	if finishedAt.Valid {
		dep.FinishedAt = &finishedAt.Time
	}
	if logs.Valid {
		dep.Logs = logs
	}
	return dep, nil

}
