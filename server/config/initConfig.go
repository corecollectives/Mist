package config

import (
	"os"

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
			log.Error().Err(err).Msg("Failed to load config file, falling back to defaults")
			Cfg = defaultConfig
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
		log.Error().Err(err).Msg("Failed to save config file, but Cfg is initialized with defaults")
		return err
	}

	log.Info().Str("path", ConfigPath).Msg("Configuration file created successfully")
	return nil
}
