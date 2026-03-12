package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
)

type RatingCacheService interface {
	InsertOrUpdateRating(userID int, animeID int, story, visuals, soundtrack uint8) error
	UpdateExistingRating(userID int, animeID int, oldStory, oldVisuals, oldSoundtrack, newStory, newVisuals, newSoundtrack uint8) error
	RemoveRating(userID int, animeID int, oldStory, oldVisuals, oldSoundtrack uint8) error
	GetRatingCache(animeID int) (*domain.RatingCache, error)
}

type RatingCacheRepository interface {
	CreateRatingCache(ctx context.Context, cache *domain.RatingCache) error
	UpdateRatingCache(ctx context.Context, cache *domain.RatingCache) error
	GetRatingCache(ctx context.Context, animeID int) (*domain.RatingCache, error)
}
