package router

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func NewRouter(db *bun.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	r := router.Group("/api/v1")
	userHandler := r.Group("/users")
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
