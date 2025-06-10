package controller

import (
	"net/http"

	"github.com/aq-simei/coin-pilot/api/service"
	responses "github.com/aq-simei/coin-pilot/internal"
	errors "github.com/aq-simei/coin-pilot/internal/config/error"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserControllerImpl struct {
	service service.UserService
}

func RegisterUserControllerRoutes(router *gin.RouterGroup, controller UserController) {
	router.GET("/:id", controller.GetUser)
	router.POST("/", controller.CreateUser)
	router.PUT("/:id", controller.UpdateUser)
	router.DELETE("/:id", controller.DeleteUser)
}

func NewUserController(service service.UserService) UserController {
	return &UserControllerImpl{
		service: service,
	}
}

func (uc *UserControllerImpl) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.service.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserControllerImpl) CreateUser(c *gin.Context) {
	var userPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		responses.BadRequest(c, "Invalid input")
		return
	}

	err := uc.service.CreateUser(c, userPayload)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case http.StatusConflict:
				responses.CustomError(c, http.StatusConflict, "User already exists")
				return
			case http.StatusInternalServerError:
				responses.InternalServerError(c, appErr.Message)
				return
			}
		}
		responses.InternalServerError(c, "Unexpected error occurred")
		return
	}

	responses.Created(c, gin.H{"message": "User created successfully"})
}

func (uc *UserControllerImpl) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var userPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := uc.service.UpdateUser(c, id, userPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (uc *UserControllerImpl) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := uc.service.DeleteUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
