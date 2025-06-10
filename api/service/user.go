package service

import (
	"context"
	"net/http"

	"github.com/aq-simei/coin-pilot/api/models"
	"github.com/aq-simei/coin-pilot/api/repository"
	errors "github.com/aq-simei/coin-pilot/internal/config/error"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (any, error)
	CreateUser(ctx context.Context, userPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}) error
	UpdateUser(ctx context.Context, id string, userPayload models.CreateUserPayload) error
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id string) (any, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		// check if it is a AppError
		if appErr, ok := err.(*errors.AppError); ok {
			if appErr.Code == http.StatusNotFound {
				return nil, errors.Wrap(404, "User not found", err)
			}

			// If it's another type of error, return it as is
			return nil, errors.Wrap(appErr.Code, appErr.Message, err)
		}
	}
	userResponse := &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, userPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
},
) error {
	// Pass errors directly from the repository
	return s.repo.CreateUser(ctx, userPayload)
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, userPayload models.CreateUserPayload,
) error {
	// Assuming the repository has an UpdateUser method
	err := s.repo.UpdateUser(ctx, id, userPayload)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	// Assuming the repository has a DeleteUser method
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
