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

func (s *UserService) UpdatePersonalInfo(ctx context.Context, id string, email *string, username *string, location *string, pronouns *string, socials *[]string) error {
	return s.userRepository.UpdatePersonalInfo(ctx, id, email, username, location, pronouns, socials)
}
