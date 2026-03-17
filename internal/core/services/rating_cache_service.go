package services

import (
	"context"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
)

type RatingCacheService struct {
	repo repositories.RatingCacheRepository
}

func NewRatingCacheService(repo repositories.RatingCacheRepository) *RatingCacheService {
	return &RatingCacheService{
		repo: repo,
	}
}

func (s *RatingCacheService) InsertOrUpdateRating(ctx context.Context, userID int, animeID int, story, visuals, soundtrack uint8) error {
	cache, err := s.repo.GetRatingCache(ctx, animeID)
	if err != nil {
		return err
	}

	if cache == nil {
		cache = domain.NewRatingCache(animeID)
		cache.UpdateCache(uint32(userID), uint32(story), uint32(visuals), uint32(soundtrack))
		return s.repo.CreateRatingCache(ctx, cache)
	}

	cache.UpdateCache(uint32(userID), uint32(story), uint32(visuals), uint32(soundtrack))
	return s.repo.UpdateRatingCache(ctx, cache)
}

func (s *RatingCacheService) UpdateExistingRating(ctx context.Context, userID int, animeID int, oldStory, oldVisuals, oldSoundtrack, newStory, newVisuals, newSoundtrack uint8) error {
	cache, err := s.repo.GetRatingCache(ctx, animeID)
	if err != nil {
		return err
	}

	if cache == nil {
		return nil
	}

	cache.UpdateExistingRating(uint32(userID), uint32(oldStory), uint32(oldVisuals), uint32(oldSoundtrack), uint32(newStory), uint32(newVisuals), uint32(newSoundtrack))
	return s.repo.UpdateRatingCache(ctx, cache)
}

func (s *RatingCacheService) RemoveRating(ctx context.Context, userID int, animeID int, oldStory, oldVisuals, oldSoundtrack uint8) error {
	cache, err := s.repo.GetRatingCache(nil, animeID)
	if err != nil {
		return err
	}

	if cache == nil {
		return nil
	}

	cache.RemoveRating(uint32(userID), uint32(oldStory), uint32(oldVisuals), uint32(oldSoundtrack))
	return s.repo.UpdateRatingCache(ctx, cache)
}

func (s *RatingCacheService) GetRatingCache(ctx context.Context, animeID int) (*domain.RatingCache, error) {
	return s.repo.GetRatingCache(ctx, animeID)
}

func (s *RatingCacheService) GetTopAnime(ctx context.Context, pageNumber, pageSize int) ([]*domain.RatingCache, utils.Pagination, error) {
	return s.GetTopAnime(ctx, pageNumber, pageSize)
}

func (s *RatingCacheService) GetPopularAnime(ctx context.Context, pageNumber, pageSize int) ([]*domain.RatingCache, utils.Pagination, error) {
	return s.GetPopularAnime(ctx, pageNumber, pageSize)
}
