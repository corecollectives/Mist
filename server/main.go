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
	// Initialize logger
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

	// make sure models get the db instance before initing the store, bcz store is dependent on models
	models.SetDB(dbInstance)
	err = store.InitStore()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing store")
		return
	}

	log.Info().Msg("Store initialized successfully")

	api.InitApiServer(dbInstance)
}
