package utils

import (
	"fmt"
	"os"
	"path/filepath"

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

	log.Info().Str("path", dynamicConfigPath).Msg("Generated Traefik dynamic config")
	return nil
}

func generateDynamicYAML(wildcardDomain *string, mistAppName string) string {
	config := `http:
  routers:`

	if wildcardDomain != nil && *wildcardDomain != "" {
		domain := *wildcardDomain

		if len(domain) > 0 && domain[0] == '*' {
			domain = domain[1:]
		}
		if len(domain) > 0 && domain[0] == '.' {
			domain = domain[1:]
		}

		mistDomain := mistAppName + "." + domain

		config += fmt.Sprintf(`
    mist-dashboard:
      rule: "Host(%s%s%s)"
      entryPoints:
        - websecure
      service: mist-dashboard
      tls:
        certResolver: le
`, "`", mistDomain, "`")

		config += fmt.Sprintf(`
    mist-dashboard-http:
      rule: "Host(%s%s%s)"
      entryPoints:
        - web
      middlewares:
        - https-redirect
      service: mist-dashboard
`, "`", mistDomain, "`")
	}

	config += `

  services:`

	if wildcardDomain != nil && *wildcardDomain != "" {
		config += `
    mist-dashboard:
      loadBalancer:
        servers:
          - url: "http://mist:5173"
`
	}

	config += `

  middlewares:
    https-redirect:
      redirectScheme:
        scheme: https
        permanent: true
`

	return config
}

func InitializeTraefikConfig(wildcardDomain *string, mistAppName string) error {
	return GenerateDynamicConfig(wildcardDomain, mistAppName)
}
