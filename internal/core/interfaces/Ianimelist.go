package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/adapters/dtos"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
)

type AnimeListRepository interface {
	AddListItem(ctx context.Context, item *domain.AnimeListItem) error
	UpdateListItem(ctx context.Context, item *domain.AnimeListItem) error
	DeleteListItem(ctx context.Context, userID int, animeID uint32) error
	FetchItemByUserAndAnime(ctx context.Context, userID int, animeID uint32) (*domain.AnimeListItem, error)
	FetchUserList(ctx context.Context, userID int, status *value.AnimeListItemStatus) ([]*domain.AnimeListItem, error)
}

type AnimeListService interface {
	AddAnimeToList(ctx context.Context, userID int, animeID uint32, status value.AnimeListItemStatus) (*dtos.AnimeListItemDTO, error)

	UpdateStatus(ctx context.Context, userID int, animeID uint32, newStatus value.AnimeListItemStatus) error
	UpdateProgress(ctx context.Context, userID int, animeID uint32, episodesWatched uint32) error
	UpdateRating(ctx context.Context, userID int, animeID uint32, story, visuals, soundtrack, enjoyment uint8) error
	UpdateNotes(ctx context.Context, userID int, animeID uint32, notes string) error
	RemoveAnimeFromList(ctx context.Context, userID int, animeID uint32) error

	FetchUserListItem(ctx context.Context, userID int, animeID uint32) (*dtos.AnimeListItemDTO, error)
	FetchUserList(ctx context.Context, userID int, status *value.AnimeListItemStatus) ([]*dtos.AnimeListItemDTO, error)
}
