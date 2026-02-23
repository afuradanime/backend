package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/utils"
)

const NotesMaxLength = 500

// Represents an entry in a user's anime list
type AnimeListItem struct {
	ID              string                    `json:"id" bson:"_id"`
	UserID          int                       `json:"userId" bson:"user_id"`
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

func NewAnimeListItem(userID int, animeID uint32, status value.AnimeListItemStatus) *AnimeListItem {
	now := time.Now()
	return &AnimeListItem{
		ID:              utils.GenerateRandomID(),
		UserID:          userID,
		AnimeID:         animeID,
		Status:          status,
		EpisodesWatched: 0,
		RewatchCount:    0,
		CreatedAt:       now,
	}
}

func (al *AnimeListItem) UpdateNotes(notes string) error {
	if len(notes) > NotesMaxLength {
		return &domain_errors.NotesLengthTooLong{
			MaxLength: NotesMaxLength,
		}
	}
	al.Notes = &notes
	return nil
}

func (al *AnimeListItem) UpdateStatus(newStatus value.AnimeListItemStatus) {
	al.Status = newStatus
	now := time.Now()
	al.EditedAt = &now
}

func (al *AnimeListItem) UpdateProgress(episodesWatched uint32, totalEpisodes uint32) error {
	// cannot watch more episodes than the anime has (unless it's still airing, in which case totalEpisodes would be 0 or nil)
	if episodesWatched > totalEpisodes && al.Status == value.AnimeListItemStatusCompleted {
		return &domain_errors.InvalidEpisodeCountErr{}
	}

	al.EpisodesWatched = episodesWatched
	now := time.Now()
	al.EditedAt = &now

	if episodesWatched == totalEpisodes {
		al.Status = value.AnimeListItemStatusCompleted
	}

	return nil
}

func (al *AnimeListItem) AddRating(story, visuals, soundtrack, enjoyment uint8) error {
	if story > 10 || visuals > 10 || soundtrack > 10 || enjoyment > 10 {
		return &domain_errors.InvalidRatingErr{}
	}

	overall := uint8((uint16(story) + uint16(visuals) + uint16(soundtrack) + uint16(enjoyment)) / 4)
	al.Rating = &Rating{
		Overall:    overall,
		Story:      story,
		Visuals:    visuals,
		Soundtrack: soundtrack,
		Enjoyment:  enjoyment,
	}

	now := time.Now()
	al.EditedAt = &now

	return nil
}

func (al *AnimeListItem) RemoveRating() {
	al.Rating = nil
	now := time.Now()
	al.EditedAt = &now
}
