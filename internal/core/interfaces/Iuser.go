package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetUserByProvider(ctx context.Context, provider string, providerID string) (*domain.User, error)
	RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdatePersonalInfo(ctx context.Context, id int, email *string, username *string, location *string, pronouns *string, socials *[]string, allowsFR, allowsRec *bool) error
	UpdateLastLogin(ctx context.Context, id int) error
}

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	GetUserByProvider(ctx context.Context, provider string, providerID string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) // same as RegisterUser
	UpdateUser(ctx context.Context, user *domain.User) error
}
