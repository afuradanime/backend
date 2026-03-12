package services

import (
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
)

type RatingCacheService struct {
	repo repositories.RatingCacheRepository
}

func NewRatingCacheService(repo repositories.RatingCacheRepository) *RatingCacheService {
	return &RatingCacheService{
		repo: repo,
	}
}

func (s *RatingCacheService) InsertOrUpdateRating(userID int, animeID int, story, visuals, soundtrack uint8) error {
	cache, err := s.repo.GetRatingCache(nil, animeID)
	if err != nil {
		return err
	}

	if cache == nil {
		cache = domain.NewRatingCache(animeID)
		cache.UpdateCache(uint32(story), uint32(visuals), uint32(soundtrack))
		return s.repo.CreateRatingCache(nil, cache)
	}

	cache.UpdateCache(uint32(story), uint32(visuals), uint32(soundtrack))
	return s.repo.UpdateRatingCache(nil, cache)
}

func (s *RatingCacheService) UpdateExistingRating(userID int, animeID int, oldStory, oldVisuals, oldSoundtrack, newStory, newVisuals, newSoundtrack uint8) error {
	cache, err := s.repo.GetRatingCache(nil, animeID)
	if err != nil {
		return err
	}

	if cache == nil {
		return nil
	}

	cache.UpdateExistingRating(uint32(oldStory), uint32(oldVisuals), uint32(oldSoundtrack), uint32(newStory), uint32(newVisuals), uint32(newSoundtrack))
	return s.repo.UpdateRatingCache(nil, cache)
}

func (s *RatingCacheService) RemoveRating(userID int, animeID int, oldStory, oldVisuals, oldSoundtrack uint8) error {
	cache, err := s.repo.GetRatingCache(nil, animeID)
	if err != nil {
		return err
	}

	if cache == nil {
		return nil
	}

	cache.RemoveRating(uint32(oldStory), uint32(oldVisuals), uint32(oldSoundtrack))
	return s.repo.UpdateRatingCache(nil, cache)
}

func (s *RatingCacheService) GetRatingCache(animeID int) (*domain.RatingCache, error) {
	return s.repo.GetRatingCache(nil, animeID)
}
