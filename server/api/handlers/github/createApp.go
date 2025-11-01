package github

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"

	"github.com/corecollectives/mist/api/middleware"
)

type HookAttributes struct {
	URL string `json:"url"`
}

type Manifest struct {
	Name               string            `json:"name"`
	URL                string            `json:"url"`
	HookAttributes     HookAttributes    `json:"hook_attributes"`
	RedirectURL        string            `json:"redirect_url"`
	SetupURL           string            `json:"setup_url"`
	CallbackURLs       []string          `json:"callback_urls"`
	Public             bool              `json:"public"`
	DefaultPermissions map[string]string `json:"default_permissions"`
	DefaultEvents      []string          `json:"default_events"`
}

func getLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no active network interface found")
}

func (h *Handler) CreateGithubApp(w http.ResponseWriter, r *http.Request) {
	appExists, err := CheckIfAppExists(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if appExists {
		http.Error(w, "GitHub App already exists", http.StatusBadRequest)
		return
	}

	userInfo, ok := middleware.GetUser(r)
	if !ok || userInfo == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userDomain := os.Getenv("VPS_DOMAIN")
	apiBase := ""
	hook := HookAttributes{}

	if userDomain != "" {
		apiBase = fmt.Sprintf("https://%s", userDomain)
		hook.URL = fmt.Sprintf("%s/api/github/webhook", apiBase)
	} else {
		ip, err := getLocalIP()
		if err != nil {
			http.Error(w, "Failed to determine server IP", http.StatusInternalServerError)
			return
		}
		apiBase = fmt.Sprintf("http://%s:8080", ip)
		hook.URL = fmt.Sprintf("%s/api/github/webhook", apiBase)
	}

	appName := fmt.Sprintf("Mist-%d", rand.Intn(90000)+10000)
	state := GenerateState(0, int(userInfo.ID))
	callbackURL := fmt.Sprintf("%s/api/github/installation/callback?state=%s", apiBase, state)

	manifest := Manifest{
		Name:           appName,
		URL:            apiBase,
		HookAttributes: hook,
		RedirectURL:    fmt.Sprintf("%s/api/github/callback", apiBase),
		SetupURL:       callbackURL,
		CallbackURLs:   []string{apiBase},
		Public:         false,
		DefaultPermissions: map[string]string{
			"contents":         "read",
			"metadata":         "read",
			"pull_requests":    "read",
			"deployments":      "read",
			"administration":   "write",
			"repository_hooks": "write",
		},
		DefaultEvents: []string{"push", "pull_request", "deployment_status"},
	}

	manifestJSON, _ := json.Marshal(manifest)
	githubManifestURL := "https://github.com/settings/apps/new?state=" + state

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<form id="manifestForm" method="post" action="%s">
			<input type="hidden" name="manifest" value='%s'/>
		</form>
		<script>document.getElementById('manifestForm').submit();</script>
	`, githubManifestURL, manifestJSON)
}
