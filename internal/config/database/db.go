package database

import (
	"context"
	"database/sql"
	"os"

	"github.com/aq-simei/coin-pilot/internal/config/environment"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewDB() *bun.DB {
	environment.LoadEnv()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		logger.Fatal("DB_DSN is not set")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// Ping to test connection
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to connect to DB: %v", err)
	}

	logger.Info("Connected to PostgreSQL")
	return db
}

func CloseDB(ctx context.Context, db *bun.DB) {
	if err := db.Close(); err != nil {
		logger.Error("Error closing DB: %v", err)
	} else {
		logger.Info("DB connection closed")
	}
}
