package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	RegisterUser(ctx context.Context, user *domain.User) error
	UpdatePersonalInfo(ctx context.Context, id string, email *string, username *string, location *string, pronouns *string, socials *[]string) error
}

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error // same as RegisterUser
	UpdatePersonalInfo(ctx context.Context, id string, email *string, username *string, location *string, pronouns *string, socials *[]string) error
}
