package services

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type RecommendationService struct {
	recommendationRepo interfaces.RecommendationRepository
	userRepo           interfaces.UserRepository
}

func NewRecommendationService(recommendationRepo interfaces.RecommendationRepository, userRepo interfaces.UserRepository) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		userRepo:           userRepo,
	}
}

func (s *RecommendationService) Send(ctx context.Context, initiatorID, receiverID, animeID int) error {

	if initiatorID == receiverID {
		return domain_errors.CannotRecommendYourselfError{}
	}

	// Check receiver allows recommendations
	receiver, err := s.userRepo.GetUserById(ctx, receiverID)
	if err != nil || receiver == nil {
		return domain_errors.UserNotFoundError{}
	}
	if !receiver.AllowsRecommendations {
		return domain_errors.RecommendationsDisabled{}
	}

	// Check stack limit
	count, err := s.recommendationRepo.RecommendationStackCount(ctx, receiverID)
	if err != nil {
		return err
	}
	if count >= domain.MAX_RECOMMENDATION_STACK {
		return domain_errors.RecommendationStackFull{}
	}

	// Check duplicate
	exists, err := s.recommendationRepo.HasBeenRecommended(ctx, receiverID, animeID)
	if err != nil {
		return err
	}
	if exists {
		return domain_errors.AlreadyRecommended{}
	}

	rec := domain.NewRecommendation(initiatorID, receiverID, animeID)
	return s.recommendationRepo.Create(ctx, rec)
}

func (s *RecommendationService) GetUserRecommendations(ctx context.Context, userID, pageNumber, pageSize int) ([]*domain.Recommendation, utils.Pagination, error) {
	return s.recommendationRepo.GetForUser(ctx, userID, pageNumber, pageSize)
}

func (s *RecommendationService) DismissRecommendation(ctx context.Context, receiverID, anime int) error {
	return s.recommendationRepo.DismissRecommendation(ctx, receiverID, anime)
}
