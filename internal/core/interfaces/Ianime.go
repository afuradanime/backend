package interfaces

import "github.com/afuradanime/backend/internal/core/domain"

type AnimeRepository interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
}

type AnimeService interface {
	FetchAnimeByID(animeID uint32) (*domain.Anime, error)
}
