package main

import (
	"fmt"

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
	api.InitApiServer(dbInstance)
}
