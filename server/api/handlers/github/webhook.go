package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/queue"
)

type WebhookPayload struct {
	Event        string              `json:"-"` // X-GitHub-Event
	Raw          json.RawMessage     `json:"-"` // full body
	Repository   *github.RepoFull    `json:"repository,omitempty"`
	Installation *github.InstallMini `json:"installation,omitempty"`
	Sender       *github.User        `json:"sender,omitempty"`
}

func (h *Handler) GithubWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ Received GitHub webhook")

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "" {
		http.Error(w, "Missing X-GitHub-Event header", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if eventType == "push" {
		var evt github.PushEvent
		if err := json.Unmarshal(body, &evt); err != nil {
			http.Error(w, "Invalid push event payload", http.StatusBadRequest)
			return
		}

		fmt.Printf("Processing push event for repo: %s\n", evt.Repository.FullName)
		depId, err := github.HandlePushEvent(evt)
		if err != nil {
			http.Error(w, "Failed to handle push event: "+err.Error(), http.StatusInternalServerError)
			return
		}
		queue := queue.GetQueue()
		queue.AddJob(depId)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}
