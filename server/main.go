package main

import (
	"fmt"
	"log"
	"time"

	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/db"

	"github.com/corecollectives/mist/store"
)

func main() {
	dbInstance, err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer dbInstance.Close()
	err = store.InitStore(dbInstance)
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return
	}
	appID := int64(1)
	commitHash := "mockcommit1234567890abcdef"
	commitMsg := "Test commit for deployment"
	triggeredBy := int64(1)
	status := "pending"
	createdAt := time.Now()

	stmt, err := dbInstance.Prepare(`
        INSERT INTO deployments(app_id, commit_hash, commit_message, triggered_by, status, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(appID, commitHash, commitMsg, triggeredBy, status, createdAt)
	if err != nil {
		log.Fatal(err)
	}
	//testing queue implementation
	api.InitApiServer(dbInstance)
}
