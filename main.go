package main

import (
	"os"

	"github.com/aq-simei/coin-pilot/api/router"
	"github.com/aq-simei/coin-pilot/internal/config/database"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load DB connection
	dbInstance := database.NewDB()
	database.RunMigrations(dbInstance)

	// Initialize Router
	router := router.NewRouter(dbInstance)

	// Read port from env or fallback
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.Info("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("failed to start server: %v", err)
	}
}
