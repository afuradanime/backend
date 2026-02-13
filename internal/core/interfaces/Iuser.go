package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	RegisterUser(ctx context.Context, user *domain.User) error
}

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error // same as RegisterUser
}
