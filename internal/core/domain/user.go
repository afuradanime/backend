package domain

import (
	"time"

	value "github.com/afuradanime/backend/internal/core/domain/value"
)

type User struct {
	ID string

	// Identity
	Email     value.Email
	Username  value.TinyStr
	AvatarURL string

	// Personal Info
	Location string
	Birthday time.Time
	Pronouns string
	Socials  []string

	// Authentication
	Provider   string
	ProviderID string

	CreatedAt time.Time
	LastLogin time.Time
}

func NewUser(id string, username string, email string) (*User, error) {
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

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}

func (u *User) UpdateAvatarURL(url string) {
	u.AvatarURL = url
}

func (u *User) UpdateLocation(location string) {
	u.Location = location
}

func (u *User) UpdateBirthday(birthday time.Time) {
	u.Birthday = birthday
}

func (u *User) UpdatePronouns(pronouns string) {
	u.Pronouns = pronouns
}

func (u *User) UpdateSocials(socials []string) {
	u.Socials = socials
}
