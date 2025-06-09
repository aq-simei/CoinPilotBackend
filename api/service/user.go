package service

import (
	"context"

	"github.com/aq-simei/coin-pilot/api/repository"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (string, error)
	CreateUser(ctx context.Context, userPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}) error
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id string) (string, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return "", err
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, userPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}) error {
	// Assuming the repository has a CreateUser method

	err := s.repo.CreateUser(ctx, userPayload)
	if err != nil {
		return err
	}
	return nil
}
