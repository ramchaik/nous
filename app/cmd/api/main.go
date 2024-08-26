package main

import (
	"log"

	"nous/internal/config"
	"nous/internal/database"
	"nous/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.New("./nous.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	srv := server.New(cfg, db)
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
