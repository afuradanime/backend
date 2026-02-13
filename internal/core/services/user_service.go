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
	// TODO: Check if user with same email or username already exists before creating a new one
	return s.userRepository.CreateUser(ctx, user)
}

func (s *UserService) UpdatePersonalInfo(ctx context.Context, id string, email *string, username *string, location *string, pronouns *string, socials *[]string) error {

	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields if new values are provided
	if email != nil {
		err := user.UpdateEmail(*email)
		if err != nil {
			return err
		}
	}

	if username != nil {
		err := user.UpdateUsername(*username)
		if err != nil {
			return err
		}
	}

	if location != nil {
		user.UpdateLocation(*location)
	}

	if pronouns != nil {
		user.UpdatePronouns(*pronouns)
	}

	if socials != nil {
		user.UpdateSocials(*socials)
	}

	return s.userRepository.UpdatePersonalInfo(ctx, id, user)
}
