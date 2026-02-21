package interfaces

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/filters"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
)

// interfaces/anime_repository.go
type AnimeRepository interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
	FetchRandomAnime() (*domain.Anime, error)
	FetchAnimeFromQuery(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
	FetchAnimeThisSeason(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
	FetchStudioByID(studioID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, utils.Pagination, error)
	FetchProducerByID(producerID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, utils.Pagination, error)
	FetchLicensorByID(licensorID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, utils.Pagination, error)
	FetchAnimeFromTag(tagID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
}

type AnimeService interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
	FetchRandomAnime() (*domain.Anime, error)
	FetchAnimeFromQuery(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
	FetchAnimeThisSeason(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)

	FetchStudioByID(filters filters.AnimeFilter, studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, utils.Pagination, error)
	FetchProducerByID(filters filters.AnimeFilter, producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, utils.Pagination, error)
	FetchLicensorByID(filters filters.AnimeFilter, licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, utils.Pagination, error)
	FetchAnimeFromTag(tagID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
}
