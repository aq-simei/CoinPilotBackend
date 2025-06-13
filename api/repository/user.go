package repository

import (
	"context"
	"net/http"

	"github.com/aq-simei/coin-pilot/api/models"
	errors "github.com/aq-simei/coin-pilot/internal/config/error"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/aq-simei/coin-pilot/internal/config/security"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, userPayload models.CreateUserPayload) error
	UpdateUser(ctx context.Context, id string, userPayload models.UpdateUserPayload) error
	DeleteUser(ctx context.Context, id string) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetUser(
	ctx context.Context,
	id string,
) (*models.User, error) {
	logger.Info("Fetching user with ID: %s", id)
	user := &models.User{}
	// Use Preload to include the Records field in the result
	result := r.db.First(user, "id = ?", id)
	if result.Error != nil {
		logger.Error("error fetching user: %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New(http.StatusNotFound, "user not found")
		}
		return nil, errors.New(http.StatusInternalServerError, "error fetching user")
	}
	return user, nil
}

func (r *UserRepositoryImpl) CreateUser(
	ctx context.Context,
	userPayload models.CreateUserPayload,
) error {
	hashedPassword, err := security.HashPassword(userPayload.Password)
	if err != nil {
		logger.Error("error hashing password: %v", err)
		return errors.New(http.StatusInternalServerError, "error hashing password")
	}
	user := &models.User{
		Name:     userPayload.Name,
		Email:    userPayload.Email,
		Password: hashedPassword,
	}

	// Check if user already exists
	existingUser := &models.User{}
	result := r.db.Where("email = ?", user.Email).First(existingUser)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		logger.Error("error checking existing user: %v", result.Error)
		return errors.New(http.StatusInternalServerError, "error checking existing user")
	}
	if result.Error == nil {
		logger.Error("attempt to create a user with an existing email: %v", user.Email)
		return errors.New(http.StatusConflict, "user already exists")
	}

	// Create the user
	if err := r.db.Create(user).Error; err != nil {
		logger.Error("error creating user: %v", err)
		return errors.New(http.StatusInternalServerError, "error creating user")
	}

	return nil
}

func (r *UserRepositoryImpl) UpdateUser(
	ctx context.Context,
	id string,
	userPayload models.UpdateUserPayload,
) error {
	// map holding non null fields
	updateData := map[string]any{}

	if userPayload.Name != nil {
		updateData["name"] = *userPayload.Name
	}
	if userPayload.Email != nil {
		updateData["email"] = *userPayload.Email
	}
	if userPayload.Password != nil {
		hashedPassword, err := security.HashPassword(*userPayload.Password)
		if err != nil {
			logger.Error("error hashing password: %v", err)
			return errors.New(http.StatusInternalServerError, "error hashing password")
		}
		updateData["password"] = hashedPassword
	}

	// Check for existing user with the same email (if email is being updated)
	if email, ok := updateData["email"]; ok {
		existingUserWithSameEmail := r.db.First(&models.User{}, "email = ? AND id != ?", email, id)
		if existingUserWithSameEmail.Error != nil && existingUserWithSameEmail.Error != gorm.ErrRecordNotFound {
			return errors.NewInternal("error checking existing user with same email")
		}
		if existingUserWithSameEmail.Error == nil {
			logger.Error("attempt to update a user with an existing email: %v", email)
			return errors.New(http.StatusConflict, "invalid email")
		}
	}

	// Perform the update
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFound("user_not_found")
	}
	return nil
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.User{})
	if result.Error != nil {
		logger.Error("error deleting user: %v", result.Error)
		return errors.New(http.StatusInternalServerError, "error deleting user")
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFound("user_not_found")
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	result := r.db.Where("email = ?", email).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New(http.StatusNotFound, "user not found")
		}
		return nil, errors.New(http.StatusInternalServerError, "error fetching user by email")
	}
	return user, nil
}
