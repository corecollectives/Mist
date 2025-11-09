package main

import (
	"fmt"

	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/db"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"

	"github.com/corecollectives/mist/store"
)

func main() {
	dbInstance, err := db.InitDB()
	_ = queue.InitQueue(dbInstance)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer dbInstance.Close()
	// make sure models get the db instance before initing the store, bcz store is dependent on models
	models.SetDB(dbInstance)
	err = store.InitStore()
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return
	}

	api.InitApiServer(dbInstance)
}
