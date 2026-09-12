package database

import (
	"log"
	"os"

	"bookend/internal/auth"
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
	if err := db.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		return nil, err
	}

	// Seed default admin user if no users exist
	if err := SeedAdminUser(db, cfg); err != nil {
		return nil, err
	}

	log.Printf("Database connection established and migrated successfully at %s", cfg.DBPath)
	return db, nil
}

// SeedAdminUser creates a default admin user if no users exist in the database
func SeedAdminUser(db *gorm.DB, cfg *config.Config) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		username := cfg.AdminUsername
		if username == "" {
			log.Fatal("No admin username set")
			os.Exit(1)
		}
		password := cfg.AdminPassword
		if password == "" {
			log.Fatal("No admin password set")
			os.Exit(1)
		}

		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			return err
		}

		admin := models.User{
			Username: username,
			Password: hashedPassword,
			Admin:    true,
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}

		log.Printf("No users found in database. Initialized default admin user '%s' (password: '%s'). Please change this password!", username, password)
	}

	return nil
}
