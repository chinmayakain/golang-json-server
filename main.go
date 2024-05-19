package main

import (
	"log"
)

func main() {
	var store *PostgresStore
	var err error

	if store, err = NewPostgresStore(); err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8080", store)
	server.Run()
}
