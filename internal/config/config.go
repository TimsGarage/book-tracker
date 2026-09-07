package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port    string
	DBPath  string
	GinMode string
}

// LoadConfig reads configuration from environment variables with fallback defaults
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "bookend.db")
	ginMode := getEnv("GIN_MODE", "debug")

	return &Config{
		Port:    port,
		DBPath:  dbPath,
		GinMode: ginMode,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
