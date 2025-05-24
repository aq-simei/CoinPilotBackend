package router

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func NewRouter(db *bun.DB) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API",
		})

	})

	return router
}
