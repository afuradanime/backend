package domain

import (
	"log"
	"slices"
	"time"

	value "github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
)

const MAX_SOCIALS = 5

type User struct {
	ID int `json:"ID" bson:"_id"`

	// Identity
	Email     value.Email    `json:"Email" bson:"email"`
	Username  value.Username `json:"Username" bson:"username"`
	AvatarURL string         `json:"AvatarURL" bson:"avatar_url"`

	// Personal Info
	Location value.TinyStr `json:"Location" bson:"location"`
	Pronouns value.TinyStr `json:"Pronouns" bson:"pronouns"`
	Birthday time.Time     `json:"Birthday" bson:"birthday"`
	Socials  []value.URL   `json:"Socials" bson:"socials"`

	// Rights
	AllowsFriendRequests  bool `json:"AllowsFriendRequests" bson:"allows_friend_requests"`
	AllowsRecommendations bool `json:"AllowsRecommendations" bson:"allows_recommendations"`
	CanPost               bool `json:"CanPost" bson:"can_post"`
	CanTranslate          bool `json:"CanTranslate" bson:"can_translate"`

	// Authentication / Authorization
	Provider   string           `json:"Provider" bson:"provider"`
	ProviderID string           `json:"ProviderID" bson:"provider_id"`
	Roles      []value.UserRole `json:"Roles" bson:"roles"`

	Badges []value.UserBadges `json:"Badges" bson:"badges"`

	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`
	LastLogin time.Time `json:"LastLogin" bson:"last_login"`
}

func NewUser(username string, email string) (*User, error) {
	newEmail, err := value.NewEmail(email)

	if err != nil {
		return nil, err
	}

	newUsername, err := value.NewUsername(username)

	if err != nil {
		return nil, err
	}

	return &User{
		// ID will be set by mongo auto-increment
		Username:              *newUsername,
		Email:                 *newEmail,
		Socials:               make([]value.URL, 0),
		Roles:                 []value.UserRole{value.UserRoleUser},
		AllowsFriendRequests:  true,
		AllowsRecommendations: true,
		CanPost:               true,
		CanTranslate:          true,
		Badges:                make([]value.UserBadges, 0),
		CreatedAt:             time.Now(),
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
	newUsername, err := value.NewUsername(username)
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

func (u *User) UpdateLocation(location string) error {
	newLocation, err := value.NewTinyStr(location)
	if err != nil {
		return err
	}

	u.Location = *newLocation
	return nil
}

func (u *User) UpdateBirthday(birthday time.Time) {
	u.Birthday = birthday
}

func (u *User) UpdatePronouns(pronouns string) error {
	newPronouns, err := value.NewTinyStr(pronouns)
	if err != nil {
		return err
	}

	u.Pronouns = *newPronouns
	return nil
}

func (u *User) UpdateSocials(socials []string) error {

	if len(socials) > MAX_SOCIALS {
		return domain_errors.TooManySocials{}
	}

	// Save to rollback if needed
	oldSocials := u.Socials

	// Assign array
	u.Socials = make([]value.URL, len(socials))

	for i := 0; i < len(socials); i++ {

		socialLink, err := value.NewURL(socials[i])
		if err != nil {
			u.Socials = oldSocials
			return err
		}

		u.Socials[i] = *socialLink
	}

	return nil
}

func (u *User) UpdateAllowsFriendRequests(allows bool) {
	u.AllowsFriendRequests = allows
}

func (u *User) UpdateAllowsRecommendations(allows bool) {
	u.AllowsRecommendations = allows
}

func (u *User) RewardBadge(badge value.UserBadges) {

	if slices.Contains(u.Badges, badge) {
		return
	}

	u.Badges = append(u.Badges, badge)
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

func (u *User) RestrictAccesses(canPost, canTranslate bool) {
	u.CanPost = canPost
	u.CanTranslate = canTranslate
}
