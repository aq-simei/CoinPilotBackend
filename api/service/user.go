package service

import (
	"context"
	"net/http"

	"github.com/aq-simei/coin-pilot/api/models"
	"github.com/aq-simei/coin-pilot/api/repository"
	errors "github.com/aq-simei/coin-pilot/internal/config/error"
	"github.com/aq-simei/coin-pilot/internal/config/security"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (any, error)
	CreateUser(ctx context.Context, userPayload models.CreateUserPayload) error
	UpdateUser(ctx context.Context, id string, userPayload models.UpdateUserPayload) error
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, claims any) error
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
		Records:   user.Records,
	}

	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, userPayload models.CreateUserPayload,
) error {
	return s.repo.CreateUser(ctx, userPayload)
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, userPayload models.UpdateUserPayload,
) error {
	err := s.repo.UpdateUser(ctx, id, userPayload)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			return "", errors.NewUnauthorized()
		}
		return "", errors.NewInternal("Failed to fetch user")
	}

	if !security.CheckPassword(user.Password, password) {
		return "", errors.NewUnauthorized()
	}

	token, err := security.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.NewInternal("Failed to generate token")
	}

	return token, nil
}

func (s *UserServiceImpl) Logout(ctx context.Context, claims any) error {
	// Implement token invalidation logic if needed (e.g., blacklist the token)
	return nil
}
