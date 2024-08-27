package main

import (
	"log"

	"nous/internal/cache"
	"nous/internal/config"
	"nous/internal/database"
	"nous/internal/llmclient"
	"nous/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	redisCache := cache.NewRedisCache(cfg.RedisAddr)

	llmClient := llmclient.NewClient(cfg.LLMBaseURL, nil, redisCache)

	srv := server.New(cfg, db, llmClient)
	if err := srv.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
