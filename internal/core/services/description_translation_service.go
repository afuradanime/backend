package services

import (
	"context"
	"slices"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type DescriptionTranslationService struct {
	translationRepository interfaces.DescriptionTranslationRepository
	animeRepository       interfaces.AnimeRepository
	userRepository        interfaces.UserRepository
}

func NewDescriptionTranslationService(
	translationRepo interfaces.DescriptionTranslationRepository,
	animeRepo interfaces.AnimeRepository,
	userRepo interfaces.UserRepository,
) *DescriptionTranslationService {
	return &DescriptionTranslationService{
		translationRepository: translationRepo,
		animeRepository:       animeRepo,
		userRepository:        userRepo,
	}
}

func (s *DescriptionTranslationService) SubmitTranslation(ctx context.Context, animeID int, translatedDescription string, createdBy int) error {

	// Check anime exists
	anime, err := s.animeRepository.FetchAnimeByID(uint32(animeID))
	if err != nil || anime == nil {
		return domain_errors.AnimeNotFoundError{AnimeID: strconv.Itoa(animeID)}
	}

	// Check no translation yet
	trans, _, _, _ := s.translationRepository.GetTranslationByAnime(ctx, animeID)
	if trans != nil {
		return domain_errors.AlreadyTranslatedError{}
	}

	// Check if user has already submitted a translation
	trans, _ = s.translationRepository.GetTranslationByAnimeFromUser(ctx, animeID, createdBy)
	if trans != nil {
		return domain_errors.AlreadySubmittedTranslation{}
	}

	// Check if user exists and can translate
	translator, err := s.userRepository.GetUserById(ctx, createdBy)
	if err != nil {
		return domain_errors.UserNotFoundError{}
	}

	if !translator.CanTranslate {
		return domain_errors.UserCantTranslate{}
	}

	// Good job
	if !slices.Contains(translator.Badges, value.UserBadgeTranslator) {
		translator.RewardBadge(value.UserBadgeTranslator)
		_ = s.userRepository.UpdateUser(ctx, translator)
	}

	translation := domain.NewDescriptionTranslation(animeID, translatedDescription, createdBy)
	return s.translationRepository.CreateTranslation(ctx, translation)
}

func (s *DescriptionTranslationService) GetMyTranslations(ctx context.Context, userID int, pageNumber, pageSize int) ([]domain.DescriptionTranslation, utils.Pagination, error) {
	return s.translationRepository.GetTranslationsByUser(ctx, userID, pageNumber, pageSize)
}

func (s *DescriptionTranslationService) GetAnimeTranslation(ctx context.Context, animeID int) (*domain.DescriptionTranslation, *domain.User, *domain.User, error) {
	t, translator, accepter, err := s.translationRepository.GetTranslationByAnime(ctx, animeID)
	if err != nil {
		return nil, nil, nil, domain_errors.TranslationNotFoundError{AnimeID: strconv.Itoa(animeID)}
	}
	return t, translator, accepter, nil
}

func (s *DescriptionTranslationService) GetPendingTranslations(ctx context.Context, pageNumber, pageSize int) ([]repositories.PendingTranslationResult, utils.Pagination, error) {
	return s.translationRepository.GetPendingTranslations(ctx, pageNumber, pageSize)
}

func (s *DescriptionTranslationService) RejectTranslation(ctx context.Context, id int, moderatorID int) error {
	t, err := s.translationRepository.GetTranslationByID(ctx, id)
	if err != nil || t == nil {
		return domain_errors.TranslationNotFoundError{}
	}

	return s.translationRepository.DeleteTranslation(ctx, id)
}

func (s *DescriptionTranslationService) AcceptTranslation(ctx context.Context, id int, moderatorID int) error {
	t, err := s.translationRepository.GetTranslationByID(ctx, id)
	if err != nil || t == nil {
		return domain_errors.TranslationNotFoundError{}
	}

	if err := t.Accept(moderatorID); err != nil {
		return err
	}

	return s.translationRepository.UpdateTranslation(ctx, t)
}
