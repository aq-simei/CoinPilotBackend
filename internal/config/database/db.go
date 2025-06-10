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
	// Use AutoMigrate for development environments
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("❌ Could not auto migrate: %v", err)
	} else {
		log.Println("✅ Auto migration ran successfully")
	}
}
