package services

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{userRepository: repo}
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) RegisterUser() error {
	return s.userRepository.CreateUser()
}
