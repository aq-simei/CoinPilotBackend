package database

import (
	"log"

	"github.com/aq-simei/coin-pilot/api/models"
	"github.com/aq-simei/coin-pilot/internal/config/environment"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	environment.LoadEnv()

	dsn := environment.GetEnv("POSTGRES_DSN")
	if dsn == "" {
		logger.Fatal("POSTGRES_DSN is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("Connected to PostgreSQL")
	return db
}

func RunMigrations(db *gorm.DB) {
	// Check if the enum type 'record_type' exists before creating it
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_type
			WHERE typname = 'record_type'
		)
	`).Scan(&exists).Error
	if err != nil {
		log.Fatalf("❌ Could not check if enum type exists: %v", err)
	}

	if !exists {
		if err := db.Exec("CREATE TYPE record_type AS ENUM ('income', 'expense')").Error; err != nil {
			log.Fatalf("❌ Could not create enum type: %v", err)
		}
		log.Println("✅ Enum type 'record_type' created successfully")
	} else {
		log.Println("ℹ️ Enum type 'record_type' already exists")
	}

	// Use AutoMigrate for development environments
	if err := db.AutoMigrate(&models.User{}, &models.Record{}); err != nil {
		log.Fatalf("❌ Could not auto migrate: %v", err)
	} else {
		log.Println("✅ Auto migration ran successfully")
	}
}
