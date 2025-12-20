package main

import (
	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/db"
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
