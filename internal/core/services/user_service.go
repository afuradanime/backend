package services

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{userRepository: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepository.GetUserById(ctx, id)
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.CreateUser(ctx, user)
}
