package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id uint) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, user *model.User) error {
	// Add business logic here (e.g., validation, password hashing)
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}
