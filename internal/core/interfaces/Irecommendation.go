package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
)

type RecommendationService interface {
	Send(ctx context.Context, initiatorID, receiverID, animeID int) error
	GetUserRecommendations(ctx context.Context, userID, pageNumber, pageSize int) ([]*domain.Recommendation, utils.Pagination, error)
	DismissRecommendation(ctx context.Context, receiverID, anime int) error
}

type RecommendationRepository interface {
	Create(ctx context.Context, rec *domain.Recommendation) error
	HasBeenRecommended(ctx context.Context, receiverID, animeID int) (bool, error)
	RecommendationStackCount(ctx context.Context, receiverID int) (int64, error)
	GetForUser(ctx context.Context, receiverID, pageNumber, pageSize int) ([]*domain.Recommendation, utils.Pagination, error)
	DismissRecommendation(ctx context.Context, receiverID, anime int) error
}
