package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

// WatchConfig watches the config file for changes and reloads it automatically.
// This function should be called as a goroutine as it runs indefinitely.
func WatchConfig() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Add the config file to watch before starting the event loop
	err = watcher.Add(ConfigPath)
	if err != nil {
		watcher.Close()
		return err
	}

	log.Info().Str("path", ConfigPath).Msg("Started watching config file for changes")

	// Run the watcher in the main goroutine (caller should run this as a goroutine)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Info().Msg("Config watcher events channel closed")
				watcher.Close()
				return nil
			}
			if event.Has(fsnotify.Write) {
				log.Info().Msg("Config file changed, reloading...")
				if err := LoadConfig(); err != nil {
					log.Error().Err(err).Msg("Failed to reload config after file change")
				} else {
					log.Info().Msg("Config reloaded successfully")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Info().Msg("Config watcher errors channel closed")
				watcher.Close()
				return nil
			}
			log.Error().Err(err).Msg("Error watching config file")
		}
	}
}
