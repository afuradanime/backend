package repositories

import (
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) CreateUser() error {
	return errors.New("not implemented")
}
