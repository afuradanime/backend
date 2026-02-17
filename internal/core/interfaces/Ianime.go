package interfaces

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
)

type AnimeRepository interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
	FetchAnimeFromQuery(name string, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
	FetchAnimeThisSeason(pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)

	FetchStudioByID(studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, utils.Pagination, error)
	FetchProducerByID(producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, utils.Pagination, error)
	FetchLicensorByID(licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, utils.Pagination, error)
}

type AnimeService interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
	FetchAnimeFromQuery(name string, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)
	FetchAnimeThisSeason(pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error)

	FetchStudioByID(studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, utils.Pagination, error)
	FetchProducerByID(producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, utils.Pagination, error)
	FetchLicensorByID(licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, utils.Pagination, error)
}
