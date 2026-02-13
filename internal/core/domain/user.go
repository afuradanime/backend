package domain

import (
	"log"
	"slices"
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

	// Authentication / Authorization
	Provider   string
	ProviderID string
	Roles      []value.UserRole

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
		Roles:    []value.UserRole{value.UserRoleUser},
	}, nil
}

func (u *User) UpdateEmail(email string) error {
	newEmail, err := value.NewEmail(email)
	if err != nil {
		return err
	}
	u.Email = *newEmail
	return nil
}

func (u *User) UpdateUsername(username string) error {
	newUsername, err := value.NewTinyStr(username)
	if err != nil {
		return err
	}
	u.Username = *newUsername
	return nil
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

func (u *User) AddRole(role value.UserRole) {

	if slices.Contains(u.Roles, role) {
		log.Printf("User %s already has role %d\n", u.Username, role)
		return
	}

	u.Roles = append(u.Roles, role)
}

func (u *User) RevokeRole(role value.UserRole) {

	if !slices.Contains(u.Roles, role) {
		log.Printf("User %s does not have role %d\n", u.Username, role)
		return
	}

	u.Roles = slices.Delete(u.Roles, slices.Index(u.Roles, role), slices.Index(u.Roles, role)+1)
}

func (u *User) HasRole(role value.UserRole) bool {
	return slices.Contains(u.Roles, role)
}
