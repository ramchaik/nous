package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	StaticPath    string
	TemplatesPath string
	ServerAddr    string
	DatabasePath  string
	LLMBaseURL    string
	RedisAddr     string
}

func Load() (*Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exPath := filepath.Dir(ex)

	return &Config{
		StaticPath:    getEnv("STATIC_PATH", filepath.Join(exPath, "..", "static")),
		TemplatesPath: getEnv("TEMPLATES_PATH", filepath.Join(exPath, "..", "templates", "*")),
		ServerAddr:    getEnv("SERVER_ADDR", ":8080"),
		DatabasePath:  getEnv("DATABASE_PATH", "./nous.db"),
		LLMBaseURL:    getEnv("LLM_BASE_URL", "http://localhost:5000"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
