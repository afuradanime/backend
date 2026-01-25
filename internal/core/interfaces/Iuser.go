package interfaces

import "github.com/afuradanime/backend/internal/core/domain"

type UserService interface {
	GetUserByID(id int) (*domain.User, error)
	RegisterUser() error // parameters TBD when auth is added
}

type UserRepository interface {
	GetUserByID(id int) (*domain.User, error)
	CreateUser() error // same as RegisterUser
}
