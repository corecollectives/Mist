package websockets

import (
	"context"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

type DeploymentEvent struct {
	Type      string      `json:"type"` // "log", "status", "progress", "error"
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type StatusUpdate struct {
	DeploymentID int64  `json:"deployment_id"`
	Status       string `json:"status"`
	Stage        string `json:"stage"`
	Progress     int    `json:"progress"`
	Message      string `json:"message"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type LogUpdate struct {
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`
}

func WatchDeploymentStatus(ctx context.Context, depID int64, events chan<- DeploymentEvent) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastStatus string
	var lastStage string
	var lastProgress int

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dep, err := models.GetDeploymentByID(depID)
			if err != nil {
				continue
			}

			if dep.Status != lastStatus || dep.Stage != lastStage || dep.Progress != lastProgress {
				lastStatus = dep.Status
				lastStage = dep.Stage
				lastProgress = dep.Progress

				errMsg := ""
				if dep.ErrorMessage != nil {
					errMsg = *dep.ErrorMessage
				}

				events <- DeploymentEvent{
					Type:      "status",
					Timestamp: time.Now(),
					Data: StatusUpdate{
						DeploymentID: depID,
						Status:       dep.Status,
						Stage:        dep.Stage,
						Progress:     dep.Progress,
						Message:      utils.GetStageMessage(dep.Stage),
						ErrorMessage: errMsg,
					},
				}
			}

			if dep.Status == "success" || dep.Status == "failed" {
				time.Sleep(500 * time.Millisecond)
				return
			}
		}
	}
}
