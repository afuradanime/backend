package domain

import (
	"log"
	"slices"
	"time"

	value "github.com/afuradanime/backend/internal/core/domain/value"
)

type User struct {
	ID int `json:"ID" bson:"_id"`

	// Identity
	Email     value.Email   `json:"Email" bson:"email"`
	Username  value.TinyStr `json:"Username" bson:"username"`
	AvatarURL string        `json:"AvatarURL" bson:"avatar_url"`

	// Personal Info
	Location string    `json:"Location" bson:"location"`
	Birthday time.Time `json:"Birthday" bson:"birthday"`
	Pronouns string    `json:"Pronouns" bson:"pronouns"`
	Socials  []string  `json:"Socials" bson:"socials"`

	// Authentication / Authorization
	Provider   string           `json:"Provider" bson:"provider"`
	ProviderID string           `json:"ProviderID" bson:"provider_id"`
	Roles      []value.UserRole `json:"Roles" bson:"roles"`

	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`
	LastLogin time.Time `json:"LastLogin" bson:"last_login"`
}

func NewUser(username string, email string) (*User, error) {
	newEmail, err := value.NewEmail(email)

	if err != nil {
		return nil, err
	}

	newUsername, err := value.NewTinyStr(username)

	if err != nil {
		return nil, err
	}

	return &User{
		// ID will be set by mongo auto-increment
		Username:  *newUsername,
		Email:     *newEmail,
		Roles:     []value.UserRole{value.UserRoleUser},
		CreatedAt: time.Now(),
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
