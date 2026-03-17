package services

import (
	"context"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type RecommendationService struct {
	recommendationRepo interfaces.RecommendationRepository
	userRepo           interfaces.UserRepository
	friendshipService  interfaces.FriendshipService
	animeListService   interfaces.AnimeListService
}

func NewRecommendationService(
	recommendationRepo interfaces.RecommendationRepository,
	userRepo interfaces.UserRepository,
	friendshipService interfaces.FriendshipService,
	animeListService interfaces.AnimeListService,
) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		userRepo:           userRepo,
		friendshipService:  friendshipService,
		animeListService:   animeListService,
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

	// Check if friends
	friendship, err := s.friendshipService.FetchFriendshipStatus(ctx, initiatorID, receiverID)
	if friendship == nil || err != nil || !friendship.AreFriends() {
		return domain_errors.NotFriendsError{}
	}

	// Check if anime is in receiver's list
	inList, err := s.animeListService.IsInAnimeList(ctx, receiverID, animeID)
	if err != nil {
		return err
	}

	if inList {
		return &domain_errors.AnimeAlreadyInListError{
			UserID:  strconv.Itoa(receiverID),
			AnimeID: strconv.Itoa(animeID),
		}
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

func (s *RecommendationService) HasBeenRecommended(ctx context.Context, receiverID, animeID int) (bool, error) {
	return s.recommendationRepo.HasBeenRecommended(ctx, receiverID, animeID)
}
