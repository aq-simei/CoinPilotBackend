package repository

import (
	"context"

	"github.com/aq-simei/coin-pilot/api/models"
	errors "github.com/aq-simei/coin-pilot/internal/config/error"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (string, error)
	CreateUser(ctx context.Context, userPayload models.CreateUserPayload) error
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
) (string, error) {
	var user string
	err := r.db.Where("id = ?", id).Find(&user)
	if err != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return "", errors.NewNotFound("User not found")
		}
	}
	return user, nil
}

func (r *UserRepositoryImpl) CreateUser(
	ctx context.Context,
	userPayload models.CreateUserPayload,
) error {
	user := &models.User{
		Name:     userPayload.Name,
		Email:    userPayload.Email,
		Password: userPayload.Password,
	}

	// check if user already exists
	existingUser := &models.User{}
	result := r.db.
		Where("email = ?", user.Email).
		First(existingUser)
	if result.Error == nil {
		return gorm.ErrDuplicatedKey
	}

	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}
