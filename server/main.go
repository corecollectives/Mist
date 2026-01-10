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
	dbInstance.AutoMigrate(
		&models.User{},
		&models.ApiToken{},
		&models.App{},
		&models.AuditLog{},
		&models.Backup{},
		&models.Deployment{},
		&models.EnvVariable{},
		&models.GithubApp{},
		&models.Project{},
		&models.ProjectMembers{},
		&models.GitProvider{},
		&models.GithubInstallation{},
		&models.AppRepositories{},
		&models.Domain{},
		&models.Volume{},
		&models.Cron{},
		&models.Registries{},
		&models.SystemSettings{},
		&models.Logs{},
		&models.AuditLog{},
		&models.ServiceTemplate{},
		&models.ApiToken{},
		&models.Session{},
		&models.Notification{},
		&models.UpdateLog{},
	)
	_ = queue.InitQueue(dbInstance)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
		return
	}
	defer dbInstance.Close()
	log.Info().Msg("Database initialized successfully")
	models.SetDB(dbInstance)

	// when we update the app, systemctl restarts the app, and we are unable to update the status of that
	// particular update in the db, and it gets stuck in 'in_progress' which leads disability in doing
	// updates, so on each startup we need to check if the last update was successfull or not and change
	// the status in the db accordingly, even if the update failed atleast we can retry it
	if err := models.CheckAndCompletePendingUpdates(); err != nil {
		log.Warn().Err(err).Msg("Failed to check pending updates")
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
