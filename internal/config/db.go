package config

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system env")
	}
}

func NewDB() *bun.DB {
	LoadEnv()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// Ping to test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}

func CloseDB(ctx context.Context, db *bun.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	} else {
		log.Println("DB connection closed")
	}
}
