package router

import (
	"github.com/aq-simei/coin-pilot/api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func NewRouter(db *bun.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	r := router.Group("/api/v1")
	userHandler := r.Group("/users").Use(middlewares.ApiKeyMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API",
		})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	userHandler.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User endpoint",
		})
	})

	return router
}
