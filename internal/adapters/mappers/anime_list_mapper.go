package mappers

import (
	"time"

	"github.com/afuradanime/backend/internal/adapters/dtos"
	"github.com/afuradanime/backend/internal/core/domain"
)

type AnimeListMapper struct{}

func NewAnimeListMapper() *AnimeListMapper {
	return &AnimeListMapper{}
}

func (m *AnimeListMapper) ToDto(li *domain.AnimeListItem, anime *domain.Anime) *dtos.AnimeListItemDTO {
	return &dtos.AnimeListItemDTO{
		AnimeID:         anime.ID,
		AnimeTitle:      anime.Title,
		AnimeEpisodes:   anime.Episodes,
		AnimeCoverURL:   anime.ImageURL,
		Status:          uint8(li.Status),
		EpisodesWatched: li.EpisodesWatched,
		Rating:          mapRatingToDto(li.Rating),
		Notes:           li.Notes,
		RewatchCount:    li.RewatchCount,
		CreatedAt:       li.CreatedAt.Format(time.RFC3339),
		EditedAt:        mapTimePtr(li.EditedAt),
	}
}

func mapRatingToDto(rating *domain.Rating) *dtos.RatingDTO {
	if rating == nil {
		return nil
	}
	return &dtos.RatingDTO{
		Overall:    rating.Overall,
		Story:      rating.Story,
		Visuals:    rating.Visuals,
		Soundtrack: rating.Soundtrack,
		Enjoyment:  rating.Enjoyment,
	}
}

func mapTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	str := t.Format(time.RFC3339)
	return &str
}
