package main

import (
	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/config"
	"github.com/corecollectives/mist/db"
	"github.com/corecollectives/mist/lib"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
	"github.com/corecollectives/mist/store"
	"github.com/corecollectives/mist/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.InitLogger()
	log.Info().Msg("Starting Mist server")
	dbInstance, err := db.InitDB()
	_ = queue.InitQueue(dbInstance)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
		return
	}
	defer dbInstance.Close()
	log.Info().Msg("Database initialized successfully")
	models.SetDB(dbInstance)

	// initialize the config file,
	// load if exists, create and then load if not exists
	if err := config.InitConfig(); err != nil {
		log.Warn().Err(err).Msg("failed to init the config file, using defaults as fallback")
	}

	// watch the .mistrc file for changes so that config can be loaded dynamically on runtime, without
	// restarting the app until there is a change in port
	// running as a background go func to not block the main thread
	go func() {
		if err := config.WatchConfig(); err != nil {
			log.Error().Err(err).Msg("Config watcher stopped with error")
		}
	}()

	// when we update the app, systemctl restarts the app, and we are unable to update the status of that
	// particular update in the db, and it gets stuck in 'in_progress' which leads disability in doing
	// updates, so on each startup we need to check if the last update was successfull or not and change
	// the status in the db accordingly, even if the update failed atleast we can retry it
	// similarly if a deployment is running and the system goes down due to overload or any other thing, it get stuck to "progress" this function cleans that too
	if err := lib.CleanupOnStartup(); err != nil {
		log.Warn().Err(err).Msg("Failed to check pending updates and deployments")
	}

	err = store.InitStore()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing store")
		return
	}
	log.Info().Msg("Store initialized successfully")
	settings, err := models.GetSystemSettings()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load system settings for Traefik initialization")
	} else {
		if err := utils.InitializeTraefikConfig(settings.WildcardDomain, settings.MistAppName); err != nil {
			log.Warn().Err(err).Msg("Failed to initialize Traefik configuration")
		} else {
			log.Info().Msg("Traefik configuration initialized successfully")
		}
	}
	api.InitApiServer()
}
