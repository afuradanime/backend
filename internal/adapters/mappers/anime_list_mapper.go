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

func (m *AnimeListMapper) ToDto(li *domain.UserListItem, anime *domain.Anime) *dtos.UserListItemDTO {
	return &dtos.UserListItemDTO{
		AnimeID:         anime.ID,
		AnimeTitle:      anime.Title,
		AnimeEpisodes:   anime.Episodes,
		AnimeCoverURL:   anime.ImageURL,
		Status:          uint8(li.Status),
		EpisodesWatched: uint32(li.EpisodesWatched),
		Rating:          mapRatingToDto(li.Rating),
		Notes:           li.Notes,
		RewatchCount:    li.RewatchCount,
		CreatedAt:       mapTime(li.CreatedAt),
		EditedAt:        mapTime(li.EditedAt),
	}
}

func mapRatingToDto(rating *uint16) *dtos.RatingDTO {
	if rating == nil {
		return nil
	}

	r := domain.Uint16ToRating(*rating)

	return &dtos.RatingDTO{
		Overall:    r.Overall,
		Story:      r.Story,
		Visuals:    r.Visuals,
		Soundtrack: r.Soundtrack,
	}
}

func mapTime(t uint32) *string {

	parsedTime := time.Unix(int64(t), 0).Format(time.RFC3339)
	return &parsedTime
}
