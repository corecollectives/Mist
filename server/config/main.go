package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
