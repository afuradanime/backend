package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
)

type DescriptionTranslationRepository interface {
	CreateTranslation(ctx context.Context, translation *domain.DescriptionTranslation) error

	GetTranslationByID(ctx context.Context, id int) (*domain.DescriptionTranslation, error)
	GetTranslationByAnime(ctx context.Context, anime int) (*domain.DescriptionTranslation, *domain.User, *domain.User, error)
	GetTranslationsByUser(ctx context.Context, userID int, pageNumber, pageSize int) ([]domain.DescriptionTranslation, utils.Pagination, error)
	GetTranslationByAnimeFromUser(ctx context.Context, anime int, id int) (*domain.DescriptionTranslation, error)

	GetPendingTranslations(ctx context.Context, pageNumber, pageSize int) ([]repositories.PendingTranslationResult, utils.Pagination, error)

	UpdateTranslation(ctx context.Context, t *domain.DescriptionTranslation) error
	DeleteTranslation(ctx context.Context, id int) error // Reject
}

type DescriptionTranslationService interface {
	SubmitTranslation(ctx context.Context, animeID int, translatedDescription string, createdBy int) error
	GetAnimeTranslation(ctx context.Context, animeID int) (*domain.DescriptionTranslation, *domain.User, *domain.User, error)

	GetPendingTranslations(ctx context.Context, pageNumber, pageSize int) ([]repositories.PendingTranslationResult, utils.Pagination, error)
	GetMyTranslations(ctx context.Context, userID int, pageNumber, pageSize int) ([]domain.DescriptionTranslation, utils.Pagination, error)

	AcceptTranslation(ctx context.Context, id int, moderatorID int) error
	RejectTranslation(ctx context.Context, id int, moderatorID int) error
}
