package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

// InitConfig initializes the config file with fully backward compatibility
// flow:
//
//  1. Check if .mistrc file exists
//     if yes - load from file and populate Cfg
//     else - go to step 2
//
//  2. Get version from database
//     if version < v1.0.3 - migrate config from database and create .mistrc file
//     else - Use default config and create .mistrc file
//
//  3. load the newly created(if previously didn't exist) config into Cfg
func InitConfig() error {
	log.Info().Msg("Initializing configuration system")

	// step 1: Check if config file exists
	if _, err := os.Stat(ConfigPath); err == nil {
		// file exists then load it and call it a day
		log.Info().Str("path", ConfigPath).Msg("Config file found, loading...")
		if err := LoadConfig(); err != nil {
			return err
		}
		log.Info().Msg("Configuration loaded successfully from file")
		return nil
	}

	// step 2: File doesnt exist then check version and create the  config file
	log.Info().Msg("Config file not found, checking version for migration")

	version, err := models.GetSystemSetting("version")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get version from database, will use default config")
		version = "1.0.0"
	}

	if version == "" {
		version = "1.0.0"
		log.Info().Msg("No version found, defaulting to 1.0.0")
	}

	log.Info().Str("version", version).Msg("Current system version")

	var newConfig ConfigType

	// compare version with v1.0.3
	if compareVersions(version, "1.0.4") < 0 {
		// version < v1.0.3 then migrate from database
		log.Info().Msg("Version < v1.0.4, migrating from database")
		migratedConfig, err := migrateFromDB()
		if err != nil {
			log.Error().Err(err).Msg("Failed to migrate from database, falling back to defaults")
			newConfig = defaultConfig
		} else {
			// merge migrated config with defaults
			newConfig = mergeWithDefaults(migratedConfig)
		}
	} else {
		// version >= v1.0.4, use default config
		log.Info().Msg("Version >= v1.0.4, using default configuration")
		newConfig = defaultConfig
	}

	// step 3: save the new config to file
	Cfg = newConfig
	if err := SaveConfig(); err != nil {
		return err
	}

	log.Info().Str("path", ConfigPath).Msg("Configuration file created successfully")
	return nil
}

// merges a config with default values for any unset fields
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

// migrateFromDB migrates config from database (for versions < v1.0.4 only)
func migrateFromDB() (ConfigType, error) {
	log.Info().Msg("Migrating config from database (version < v1.0.4)")

	config := ConfigType{}

	// get system settings from database
	settings, err := models.GetSystemSettings()
	if err != nil {
		return config, fmt.Errorf("failed to get system settings: %w", err)
	}

	// map database settings to config structure
	// network settings
	config.Network.WildcardDomain = settings.WildcardDomain
	config.Network.MistAppName = &settings.MistAppName

	// security settings
	config.Security.SecureCookies = &settings.SecureCookies

	// docker settings
	config.Docker.AutoCleanupContainers = &settings.AutoCleanupContainers

	log.Info().Msg("Successfully migrated config from database")
	return config, nil
}

// returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
// i think we can directly compare the strings, might change this later
func compareVersions(v1, v2 string) int {
	// remove 'v' prefix
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
