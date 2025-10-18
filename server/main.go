package main

import (
	"fmt"

	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/db"
)

func main() {
	dbInstance, err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer dbInstance.Close()
	fmt.Println("Database initialized")
	api.InitApiServer()
}
