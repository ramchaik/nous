package config

import (
	"os"
	"path/filepath"
)

// Config holds all configuration for the application
type Config struct {
    StaticPath    string
    TemplatesPath string
    ServerAddr    string
    DatabasePath  string
    LLMBaseURL    string
    RedisAddr     string
}

// Load returns a new Config struct populated with values from environment variables or defaults
func Load() (*Config, error) {
    exPath, err := getExecutablePath()
    if err != nil {
        return nil, err
    }

    return &Config{
        StaticPath:    getEnv("STATIC_PATH", filepath.Join(exPath, "..", "static")),
        TemplatesPath: getEnv("TEMPLATES_PATH", filepath.Join(exPath, "..", "templates", "*")),
        ServerAddr:    getEnv("SERVER_ADDR", ":8080"),
        DatabasePath:  getEnv("DATABASE_PATH", "./nous.db"),
        LLMBaseURL:    getEnv("LLM_BASE_URL", "http://localhost:5000"),
        RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
    }, nil
}

// getExecutablePath returns the directory of the current executable
func getExecutablePath() (string, error) {
    ex, err := os.Executable()
    if err != nil {
        return "", err
    }
    return filepath.Dir(ex), nil
}

// getEnv retrieves the value of the environment variable named by the key
// If the variable is not present, it returns the fallback value
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
