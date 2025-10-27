package store

import (
	"database/sql"
	"sync"
)

// var SetupRequired bool = true
type SetupState struct {
	setupRequired bool
	mu            sync.RWMutex
}

var state = &SetupState{setupRequired: true}

func InitSetupRequired(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return err
	}
	state.mu.Lock()
	state.setupRequired = count == 0
	state.mu.Unlock()
	return nil

}

func SetSetupRequired(setupRequired bool) {
	state.mu.Lock()
	state.setupRequired = setupRequired
	state.mu.Unlock()

}

func IsSetupRequired() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.setupRequired
}
