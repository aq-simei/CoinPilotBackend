package router

import (
	"github.com/aq-simei/coin-pilot/api/controller"
	"github.com/aq-simei/coin-pilot/api/middlewares"
	"github.com/aq-simei/coin-pilot/api/repository"
	"github.com/aq-simei/coin-pilot/api/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
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
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	userHandler.GET("/", userController.GetUser)
	userHandler.POST("/", userController.CreateUser)
	userHandler.PUT("/:id", userController.UpdateUser)

	return router
}
