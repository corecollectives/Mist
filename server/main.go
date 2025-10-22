package main

import (
	"fmt"
	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/db"
)

func main() {
	dbInstance, err := db.InitDB()
	handlers.SetDB(dbInstance)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer dbInstance.Close()
	api.InitApiServer()
}
