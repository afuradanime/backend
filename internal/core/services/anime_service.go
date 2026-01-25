package services

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type AnimeService struct {
	repo interfaces.AnimeRepository
}

func NewAnimeService(repo interfaces.AnimeRepository) *AnimeService {
	return &AnimeService{repo: repo}
}

func (s *AnimeService) FetchAnimeByID(animeID uint32) (*domain.Anime, error) {
	return s.repo.FetchAnimeByID(animeID)
}
