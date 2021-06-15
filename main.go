package main

import (
	"log"

	"github.com/mpfen/Go-Todo-REST-API-V2/api"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
)

func main() {
	db := store.NewDatabaseConnection("database.db")
	server := api.NewTodoServer(db)

	err := server.Router.Run(":5000")

	if err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
