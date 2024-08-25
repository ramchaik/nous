package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	StaticPath    string
	TemplatesPath string
	Port          string
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
		Port:          ":8080",
	}, nil
}
