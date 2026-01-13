package config

import "github.com/moby/moby/api/types/container"

var defaultConfig = ConfigType{
	Server: ServerConfig{
		Port:                 ptr(8080),
		DefaultAppPort:       ptr(3000),
		APIReadHeaderTimeout: ptr(10), // seconds
		MaxAvatarSize:        ptr(5),  // MB
	},
	Network: NetworkConfig{
		DNSValidationTimeout: ptr(60),
	},
	Security: SecurityConfig{
		JWTExpiry:         ptr(744), // hours
		SecureCookies:     ptr(true),
		PasswordMinLength: ptr(8),
	},
	Docker: DockerConfig{
		AutoCleanupContainers:   ptr(true),
		DefaultRestartPolicy:    ptr(container.RestartPolicyAlways),
		BuildImageTimeout:       ptr(10),
		PullImageTimeout:        ptr(5),
		StartContainerTimeout:   ptr(1),
		StopContainerTimeout:    ptr(1),
		RestartContainerTimeout: ptr(1),
	},
	Git: GitConfig{
		GitCloneTimeout:         ptr(5),
		RemoveGitRepoAfterBuild: ptr(true),
	},
}

func ptr[T any](v T) *T {
	return &v
}

func getOr[T any](v *T, fallback *T) T {
	if v != nil {
		return *v
	}
	return *fallback
}
