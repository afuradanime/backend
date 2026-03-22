package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
)

type Group struct {
	ID          int           `json:"ID" bson:"_id"`
	Name        value.TinyStr `json:"Name" bson:"name"`
	Icon        value.URL     `json:"Icon" bson:"icon"`
	Description value.LongStr `json:"Description" bson:"description"`
	Rules       value.LongStr `json:"Rules" bson:"rules"`
	Public      bool          `json:"Public" bson:"public"`

	Moderators []int `json:"Mods" bson:"mods"`

	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`
}

func NewGroup(name, description, rules, icon string) (*Group, error) {
	nameVal, err := value.NewTinyStr(name)
	if err != nil {
		return nil, err
	}

	descVal, err := value.NewLongStr(description)
	if err != nil {
		return nil, err
	}

	rulesVal, err := value.NewLongStr(rules)
	if err != nil {
		return nil, err
	}

	iconVal, err := value.NewURL(icon)
	if err != nil {
		return nil, err
	}

	return &Group{
		Name:        *nameVal,
		Icon:        *iconVal,
		Description: *descVal,
		Rules:       *rulesVal,
		Moderators:  make([]int, 0),
		Public:      true,
		CreatedAt:   time.Now(),
	}, nil
}

func (g *Group) UpdateName(name string) error {

	newName, err := value.NewTinyStr(name)
	if err != nil {
		return err
	}

	g.Name = *newName
	return nil
}

func (g *Group) UpdateDescription(description string) error {

	desc, err := value.NewLongStr(description)
	if err != nil {
		return err
	}

	g.Description = *desc
	return nil
}

func (g *Group) UpdateRules(rules string) error {

	r, err := value.NewLongStr(rules)
	if err != nil {
		return err
	}

	g.Rules = *r
	return nil
}

func (g *Group) UpdateIcon(icon string) error {

	i, err := value.NewURL(icon)
	if err != nil {
		return err
	}

	g.Icon = *i
	return nil
}

func (g *Group) AddModerator(userID int) error {

	// check for duplicates
	for _, id := range g.Moderators {
		if id == userID {
			return domain_errors.AlreadyModeratingError{}
		}
	}

	g.Moderators = append(g.Moderators, userID)
	return nil
}

func (g *Group) RemoveModerator(userID int) error {

	if len(g.Moderators) == 1 {
		return domain_errors.NoModeratorsLeftError{}
	}

	index := -1
	for i, id := range g.Moderators {
		if id == userID {
			index = i
			break
		}
	}

	if index == -1 {
		return domain_errors.NotModeratingError{}
	}

	g.Moderators[index] = g.Moderators[len(g.Moderators)-1]
	g.Moderators = g.Moderators[:len(g.Moderators)-1]

	return nil
}

func (g *Group) IsModerator(userID int) bool {
	for _, id := range g.Moderators {
		if id == userID {
			return true
		}
	}
	return false
}

func (g *Group) MakePrivate() {
	g.Public = false
}
