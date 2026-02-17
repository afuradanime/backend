package interfaces

import (
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
)

type AnimeRepository interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
	FetchAnimeFromQuery(name string, pageNumber, pageSize int) ([]*domain.Anime, error)
	FetchAnimeThisSeason() ([]*domain.Anime, error)

	FetchStudioByID(studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, error)
	FetchProducerByID(producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, error)
	FetchLicensorByID(licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, error)
}

type AnimeService interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
	FetchAnimeFromQuery(name string, pageNumber, pageSize int) ([]*domain.Anime, error)
	FetchAnimeThisSeason() ([]*domain.Anime, error)

	FetchStudioByID(studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, error)
	FetchProducerByID(producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, error)
	FetchLicensorByID(licenserID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, error)
}
