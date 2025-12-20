package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	TraefikConfigDir   = "/var/lib/mist/traefik"
	TraefikStaticFile  = "traefik.yml"
	TraefikDynamicFile = "dynamic.yml"
)

func GenerateDynamicConfig(wildcardDomain *string, mistAppName string) error {
	if err := os.MkdirAll(TraefikConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create traefik config directory: %w", err)
	}

	dynamicConfigPath := filepath.Join(TraefikConfigDir, TraefikDynamicFile)

	content := generateDynamicYAML(wildcardDomain, mistAppName)

	if err := os.WriteFile(dynamicConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write dynamic config: %w", err)
	}

	log.Info().
		Str("path", dynamicConfigPath).
		Msg("Generated Traefik dynamic config")

	return nil
}

func generateDynamicYAML(wildcardDomain *string, mistAppName string) string {
	var b strings.Builder

	b.WriteString(`http:
  routers:
`)

	var mistDomain string
	if wildcardDomain != nil && *wildcardDomain != "" {
		domain := strings.TrimPrefix(*wildcardDomain, "*")
		domain = strings.TrimPrefix(domain, ".")

		mistDomain = mistAppName + "." + domain

		b.WriteString(fmt.Sprintf(`
    mist-dashboard:
      rule: "Host(`+"`%s`"+`)"
      entryPoints:
        - websecure
      service: mist-dashboard
      tls:
        certResolver: le

    mist-dashboard-http:
      rule: "Host(`+"`%s`"+`)"
      entryPoints:
        - web
      middlewares:
        - https-redirect
      service: mist-dashboard
`, mistDomain, mistDomain))
	}

	b.WriteString(`
  services:
`)

	if mistDomain != "" {
		b.WriteString(`
    mist-dashboard:
      loadBalancer:
        servers:
          - url: "http://172.17.0.1:8080"
`)
	}

	b.WriteString(`
  middlewares:
    https-redirect:
      redirectScheme:
        scheme: https
        permanent: true
`)

	return b.String()
}

func InitializeTraefikConfig(wildcardDomain *string, mistAppName string) error {
	return GenerateDynamicConfig(wildcardDomain, mistAppName)
}
