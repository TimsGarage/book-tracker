package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port          string
	DBPath        string
	GinMode       string
	JWTSecret     string
	AdminUsername string
	AdminPassword string
}

// LoadConfig reads configuration from environment variables with fallback defaults
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "bookend.db")
	ginMode := getEnv("GIN_MODE", "debug")
	jwtSecret := getEnv("JWT_SECRET", "bookend-default-secret-change-me")
	adminUsername := getEnv("ADMIN_USERNAME", "admin")
	adminPassword := getEnv("ADMIN_PASSWORD", "admin123")

	return &Config{
		Port:          port,
		DBPath:        dbPath,
		GinMode:       ginMode,
		JWTSecret:     jwtSecret,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
