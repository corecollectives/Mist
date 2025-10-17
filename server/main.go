package main

import (
	"fmt"

	"github.com/corecollectives/mist/api"
	"github.com/corecollectives/mist/db"
)

func main() {
	dbInstance, err := db.InitDB("mist.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database initialized:", dbInstance)
	api.InitApiServer()
}
