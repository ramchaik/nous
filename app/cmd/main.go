package main

import (
	"log"

	"nous/internal/config"
	"nous/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	s := server.New(cfg)
	if err := s.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
