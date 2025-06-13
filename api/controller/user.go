package controller

import (
	"net/http"

	"github.com/aq-simei/coin-pilot/api/models"
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
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type UserControllerImpl struct {
	service service.UserService
}

func RegisterUserControllerRoutes(router *gin.RouterGroup, controller UserController) {
	router.GET("/:id", controller.GetUser)
	router.POST("/", controller.CreateUser)
	router.PUT("/:id", controller.UpdateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.POST("/login", controller.Login)
	router.POST("/logout", controller.Logout)
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
		responses.InternalServerError(c, "Failed to retrieve user")
		return
	}
	responses.Success(c, user)
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
	var userPayload models.UpdateUserPayload
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		responses.BadRequest(c, "Invalid input")
		return
	}
	err := uc.service.UpdateUser(c, id, userPayload)
	if err != nil {
		responses.InternalServerError(c, "Failed to update user")
		return
	}
	responses.Success(c, "User updated successfully")
}

func (uc *UserControllerImpl) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := uc.service.DeleteUser(c, id)
	if err != nil {
		responses.InternalServerError(c, "Failed to delete user")
		return
	}
	responses.Success(c, "User deleted successfully")
}

func (uc *UserControllerImpl) Login(c *gin.Context) {
	var loginPayload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		responses.BadRequest(c, "Invalid input")
		return
	}

	token, err := uc.service.Login(c, loginPayload.Email, loginPayload.Password)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case http.StatusUnauthorized:
				responses.Unauthorized(c, appErr.Message)
				return
			case http.StatusInternalServerError:
				responses.InternalServerError(c, appErr.Message)
				return
			}
		}
		responses.InternalServerError(c, "Unexpected error occurred")
		return
	}

	responses.Success(c, gin.H{"token": token})
}

func (uc *UserControllerImpl) Logout(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		responses.Unauthorized(c, "Invalid token")
		return
	}

	err := uc.service.Logout(c, claims)
	if err != nil {
		responses.InternalServerError(c, "Failed to logout")
		return
	}

	responses.Success(c, "User logged out successfully")
}
