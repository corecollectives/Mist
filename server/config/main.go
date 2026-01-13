package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

const ConfigPath = "/var/lib/mist/.mistrc"

var Cfg ConfigType

func WriteConfig(config ConfigType) {
	Cfg = config
}

func SaveConfig() error {
	dir := filepath.Dir(ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func LoadConfig() error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found at %s", ConfigPath)
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

type ResolvedConfig struct {
	Server   ResolvedServerConfig
	Network  ResolvedNetworkConfig
	Security ResolvedSecurityConfig
	Docker   ResolvedDockerConfig
	Git      ResolvedGitConfig
}

func GetConfig() ResolvedConfig {
	return ResolvedConfig{
		Server: ResolvedServerConfig{
			Port:                 getOr(Cfg.Server.Port, defaultConfig.Server.Port),
			DefaultAppPort:       getOr(Cfg.Server.DefaultAppPort, defaultConfig.Server.DefaultAppPort),
			APIReadHeaderTimeout: getOr(Cfg.Server.APIReadHeaderTimeout, defaultConfig.Server.APIReadHeaderTimeout),
			MaxAvatarSize:        getOr(Cfg.Server.MaxAvatarSize, defaultConfig.Server.MaxAvatarSize),
		},
		Network: ResolvedNetworkConfig{
			WildcardDomain:       getOr(Cfg.Network.WildcardDomain, defaultConfig.Network.WildcardDomain),
			MistAppName:          getOr(Cfg.Network.MistAppName, defaultConfig.Network.MistAppName),
			DNSValidationTimeout: getOr(Cfg.Network.DNSValidationTimeout, defaultConfig.Network.DNSValidationTimeout),
		},
		Security: ResolvedSecurityConfig{
			JWTExpiry:         getOr(Cfg.Security.JWTExpiry, defaultConfig.Security.JWTExpiry),
			SecureCookies:     getOr(Cfg.Security.SecureCookies, defaultConfig.Security.SecureCookies),
			PasswordMinLength: getOr(Cfg.Security.PasswordMinLength, defaultConfig.Security.PasswordMinLength),
		},
		Docker: ResolvedDockerConfig{
			AutoCleanupContainers:   getOr(Cfg.Docker.AutoCleanupContainers, defaultConfig.Docker.AutoCleanupContainers),
			DefaultRestartPolicy:    getOr(Cfg.Docker.DefaultRestartPolicy, defaultConfig.Docker.DefaultRestartPolicy),
			BuildImageTimeout:       getOr(Cfg.Docker.BuildImageTimeout, defaultConfig.Docker.BuildImageTimeout),
			PullImageTimeout:        getOr(Cfg.Docker.PullImageTimeout, defaultConfig.Docker.PullImageTimeout),
			StartContainerTimeout:   getOr(Cfg.Docker.StartContainerTimeout, defaultConfig.Docker.StartContainerTimeout),
			StopContainerTimeout:    getOr(Cfg.Docker.StopContainerTimeout, defaultConfig.Docker.StopContainerTimeout),
			RestartContainerTimeout: getOr(Cfg.Docker.RestartContainerTimeout, defaultConfig.Docker.RestartContainerTimeout),
		},
		Git: ResolvedGitConfig{
			GitCloneTimeout:         getOr(Cfg.Git.GitCloneTimeout, defaultConfig.Git.GitCloneTimeout),
			RemoveGitRepoAfterBuild: getOr(Cfg.Git.RemoveGitRepoAfterBuild, defaultConfig.Git.RemoveGitRepoAfterBuild),
		},
	}
}

// compareVersions compares two semantic versions (e.g., "1.0.2" vs "1.0.3")
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var num1, num2 int
		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		if num1 < num2 {
			return -1
		}
		if num1 > num2 {
			return 1
		}
	}

	return 0
}

// migrateFromDB migrates config from database (for versions < v1.0.4)
func migrateFromDB() (ConfigType, error) {
	log.Info().Msg("Migrating config from database (version < v1.0.4)")

	config := ConfigType{}

	settings, err := models.GetSystemSettings()
	if err != nil {
		return config, fmt.Errorf("failed to get system settings: %w", err)
	}

	config.Network.WildcardDomain = settings.WildcardDomain
	config.Network.MistAppName = &settings.MistAppName
	config.Security.SecureCookies = &settings.SecureCookies
	config.Docker.AutoCleanupContainers = &settings.AutoCleanupContainers

	log.Info().Msg("Successfully migrated config from database")
	return config, nil
}

// mergeWithDefaults merges a config with default values for any unset fields
func mergeWithDefaults(config ConfigType) ConfigType {
	merged := config

	if merged.Server.Port == nil {
		merged.Server.Port = defaultConfig.Server.Port
	}
	if merged.Server.DefaultAppPort == nil {
		merged.Server.DefaultAppPort = defaultConfig.Server.DefaultAppPort
	}
	if merged.Server.APIReadHeaderTimeout == nil {
		merged.Server.APIReadHeaderTimeout = defaultConfig.Server.APIReadHeaderTimeout
	}
	if merged.Server.MaxAvatarSize == nil {
		merged.Server.MaxAvatarSize = defaultConfig.Server.MaxAvatarSize
	}
	if merged.Network.MistAppName == nil {
		merged.Network.MistAppName = defaultConfig.Network.MistAppName
	}
	if merged.Network.DNSValidationTimeout == nil {
		merged.Network.DNSValidationTimeout = defaultConfig.Network.DNSValidationTimeout
	}
	if merged.Security.JWTExpiry == nil {
		merged.Security.JWTExpiry = defaultConfig.Security.JWTExpiry
	}
	if merged.Security.SecureCookies == nil {
		merged.Security.SecureCookies = defaultConfig.Security.SecureCookies
	}
	if merged.Security.PasswordMinLength == nil {
		merged.Security.PasswordMinLength = defaultConfig.Security.PasswordMinLength
	}
	if merged.Docker.AutoCleanupContainers == nil {
		merged.Docker.AutoCleanupContainers = defaultConfig.Docker.AutoCleanupContainers
	}
	if merged.Docker.DefaultRestartPolicy == nil {
		merged.Docker.DefaultRestartPolicy = defaultConfig.Docker.DefaultRestartPolicy
	}
	if merged.Docker.BuildImageTimeout == nil {
		merged.Docker.BuildImageTimeout = defaultConfig.Docker.BuildImageTimeout
	}
	if merged.Docker.PullImageTimeout == nil {
		merged.Docker.PullImageTimeout = defaultConfig.Docker.PullImageTimeout
	}
	if merged.Docker.StartContainerTimeout == nil {
		merged.Docker.StartContainerTimeout = defaultConfig.Docker.StartContainerTimeout
	}
	if merged.Docker.StopContainerTimeout == nil {
		merged.Docker.StopContainerTimeout = defaultConfig.Docker.StopContainerTimeout
	}
	if merged.Docker.RestartContainerTimeout == nil {
		merged.Docker.RestartContainerTimeout = defaultConfig.Docker.RestartContainerTimeout
	}
	if merged.Git.GitCloneTimeout == nil {
		merged.Git.GitCloneTimeout = defaultConfig.Git.GitCloneTimeout
	}
	if merged.Git.RemoveGitRepoAfterBuild == nil {
		merged.Git.RemoveGitRepoAfterBuild = defaultConfig.Git.RemoveGitRepoAfterBuild
	}

	return merged
}
