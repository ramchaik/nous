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
}

func Load() (*Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exPath := filepath.Dir(ex)

	return &Config{
		StaticPath:    filepath.Join(exPath, "..", "static"),
		TemplatesPath: filepath.Join(exPath, "..", "templates", "*"),
		ServerAddr:    ":8080",
		DatabasePath:  "./nous.db",
		LLMBaseURL:    "http://localhost:5000",
	}, nil
}
