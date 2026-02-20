package services

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/filters"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
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

func (s *AnimeService) FetchAnimeFromQuery(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error) {
	return s.repo.FetchAnimeFromQuery(filters, pageNumber, pageSize)
}

func (s *AnimeService) FetchAnimeThisSeason(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error) {
	return s.repo.FetchAnimeThisSeason(filters, pageNumber, pageSize)
}

func (s *AnimeService) FetchStudioByID(filters filters.AnimeFilter, studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, utils.Pagination, error) {
	return s.repo.FetchStudioByID(studioID, filters, pageNumber, pageSize)
}

func (s *AnimeService) FetchProducerByID(filters filters.AnimeFilter, producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, utils.Pagination, error) {
	return s.repo.FetchProducerByID(producerID, filters, pageNumber, pageSize)
}

func (s *AnimeService) FetchLicensorByID(filters filters.AnimeFilter, licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, utils.Pagination, error) {
	return s.repo.FetchLicensorByID(licenserID, filters, pageNumber, pageSize)
}

func (s *AnimeService) FetchAnimeFromTag(tagID uint32, f filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error) {
	return s.repo.FetchAnimeFromTag(tagID, f, pageNumber, pageSize)
}
