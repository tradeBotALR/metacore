package main

import (
	"log"
	"metacore/configs"
	"metacore/postgres"
)

func main() {
	cfg := configs.DefaultConfig()

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
}
