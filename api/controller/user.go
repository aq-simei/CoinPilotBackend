package controller

import (
	"net/http"

	"github.com/aq-simei/coin-pilot/api/service"
	responses "github.com/aq-simei/coin-pilot/internal"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController interface {
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userController{
		service: service,
	}
}

func (uc *userController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.service.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *userController) CreateUser(c *gin.Context) {
	var userPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// Assuming the service has a CreateUser method
	err := uc.service.CreateUser(c, userPayload)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			responses.InternalServerError(c, "User already exists")
			return
		}
		return
	}
}
