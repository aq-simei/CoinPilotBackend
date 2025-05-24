package main

import (
	"context"
	"log"
	"os"

	"github.com/aq-simei/coin-pilot/api/router"
	"github.com/aq-simei/coin-pilot/internal/config"
)

func main() {
	// Load DB connection

	dbInstance := config.NewDB()
	defer config.CloseDB(context.Background(), dbInstance)

	// Initialize Router
	router := router.NewRouter(dbInstance)

	// Read port from env or fallback
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
