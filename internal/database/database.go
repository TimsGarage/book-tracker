package database

import (
	"log"

	"bookend/internal/config"
	"bookend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the GORM database connection and migrates the schemas
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Info
	if cfg.GinMode == "release" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Running database auto-migrations...")
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		return nil, err
	}

	log.Printf("Database connection established and migrated successfully at %s", cfg.DBPath)
	return db, nil
}
