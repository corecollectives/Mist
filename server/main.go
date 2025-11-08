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
	err = store.InitStore(dbInstance)
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return
	}
	models.SetDB(dbInstance)

	api.InitApiServer(dbInstance)
}
