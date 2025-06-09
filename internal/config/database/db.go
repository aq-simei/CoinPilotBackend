package database

import (
	"log"

	"github.com/aq-simei/coin-pilot/api/models"
	"github.com/aq-simei/coin-pilot/internal/config/environment"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/go-gormigrate/gormigrate/v2"
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
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20240608_create_users",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateTable(&models.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("❌ Could not migrate: %v", err)
	} else {
		log.Println("✅ Migrations ran successfully")
	}
}
