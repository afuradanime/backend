package domain

import (
	"time"

	value "github.com/afuradanime/backend/internal/core/domain/value"
)

type User struct {
	ID int

	Email     value.Email
	Username  value.TinyStr
	AvatarURL string

	Provider   string
	ProviderID string

	CreatedAt time.Time
	LastLogin time.Time
}

func NewUser(id int, username string, email string) (*User, error) {
	newEmail, err := value.NewEmail(email)

	if err != nil {
		return nil, err
	}

	newUsername, err := value.NewTinyStr(username)

	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Username: *newUsername,
		Email:    *newEmail,
	}, nil
}
