package services

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
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

func (s *AnimeService) FetchAnimeFromQuery(name string, pageNumber, pageSize int) ([]*domain.Anime, error) {
	return s.repo.FetchAnimeFromQuery(name, pageNumber, pageSize)
}

func (s *AnimeService) FetchAnimeThisSeason() ([]*domain.Anime, error) {
	return s.repo.FetchAnimeThisSeason()
}

func (s *AnimeService) FetchStudioByID(studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, error) {
	return s.repo.FetchStudioByID(studioID, pageNumber, pageSize)
}

func (s *AnimeService) FetchProducerByID(producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, error) {
	return s.repo.FetchProducerByID(producerID, pageNumber, pageSize)
}

func (s *AnimeService) FetchLicensorByID(licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, error) {
	return s.repo.FetchLicensorByID(licenserID, pageNumber, pageSize)
}
