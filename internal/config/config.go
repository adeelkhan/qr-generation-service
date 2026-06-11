package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load reads configuration from environment variables (and an optional .env file).
func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        port,
	}
}
