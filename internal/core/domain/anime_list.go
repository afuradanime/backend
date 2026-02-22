package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
)

// Intermediary struct representing a user's anime list, unique for each user
type AnimeList struct {
	UserID int              `json:"userId" bson:"user_id"`
	Items  []*AnimeListItem `json:"items" bson:"items"`
}

// Represents an entry in a user's anime list
type AnimeListItem struct {
	ID              string                    `json:"id" bson:"_id"`
	AnimeID         uint32                    `json:"animeId" bson:"anime_id"`
	Status          value.AnimeListItemStatus `json:"status" bson:"status"`
	EpisodesWatched uint32                    `json:"episodesWatched" bson:"episodes_watched"`
	Rating          *Rating                   `json:"rating,omitempty" bson:"rating,omitempty"`
	Notes           *string                   `json:"notes,omitempty" bson:"notes,omitempty"`
	RewatchCount    uint8                     `json:"rewatchCount" bson:"rewatch_count"`
	CreatedAt       time.Time                 `json:"createdAt" bson:"created_at"`
	EditedAt        *time.Time                `json:"editedAt,omitempty" bson:"edited_at,omitempty"`
}

// Represents the rating a user gives to an anime in their list, with an overall rating and optional breakdown by category
type Rating struct {
	Overall    uint8 `json:"overall" bson:"overall"`
	Story      uint8 `json:"story" bson:"story"`
	Visuals    uint8 `json:"visuals" bson:"visuals"`
	Soundtrack uint8 `json:"soundtrack" bson:"soundtrack"`
	Enjoyment  uint8 `json:"enjoyment" bson:"enjoyment"`
}

func NewAnimeList(userID int) *AnimeList {
	return &AnimeList{
		UserID: userID,
		Items:  []*AnimeListItem{},
	}
}

func NewAnimeListItem(animeID uint32, status value.AnimeListItemStatus, episodesWatched uint32, rating *Rating, notes *string) *AnimeListItem {
	now := time.Now()
	return &AnimeListItem{
		ID:              utils.GenerateRandomID(),
		AnimeID:         animeID,
		Status:          status,
		EpisodesWatched: episodesWatched,
		Rating:          rating,
		Notes:           notes,
		RewatchCount:    0,
		CreatedAt:       now,
		EditedAt:        nil,
	}
}

func NewRating(story, visuals, soundtrack, enjoyment uint8) *Rating {
	overall := uint8((uint16(story) + uint16(visuals) + uint16(soundtrack) + uint16(enjoyment)) / 4)
	return &Rating{
		Overall:    overall,
		Story:      story,
		Visuals:    visuals,
		Soundtrack: soundtrack,
		Enjoyment:  enjoyment,
	}
}
