package services

import (
	"context"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type DescriptionTranslationService struct {
	translationRepository interfaces.DescriptionTranslationRepository
	animeRepository       interfaces.AnimeRepository
}

func NewDescriptionTranslationService(
	translationRepo interfaces.DescriptionTranslationRepository,
	animeRepo interfaces.AnimeRepository,
) *DescriptionTranslationService {
	return &DescriptionTranslationService{
		translationRepository: translationRepo,
		animeRepository:       animeRepo,
	}
}

func (s *DescriptionTranslationService) SubmitTranslation(ctx context.Context, animeID int, translatedDescription string, createdBy int) error {

	// Check anime exists
	anime, err := s.animeRepository.FetchAnimeByID(uint32(animeID))
	if err != nil || anime == nil {
		return domain_errors.AnimeNotFoundError{AnimeID: strconv.Itoa(animeID)}
	}

	// Check no translation yet
	trans, _ := s.translationRepository.GetTranslationByAnime(ctx, animeID)
	if trans != nil {
		return domain_errors.AlreadyTranslatedError{}
	}

	// Check if user has already submitted a translation
	trans, _ = s.translationRepository.GetTranslationByAnimeFromUser(ctx, animeID, createdBy)
	if trans != nil {
		return domain_errors.AlreadySubmittedTranslation{}
	}

	translation := domain.NewDescriptionTranslation(animeID, translatedDescription, createdBy)
	return s.translationRepository.CreateTranslation(ctx, translation)
}

func (s *DescriptionTranslationService) GetMyTranslations(ctx context.Context, userID int, pageNumber, pageSize int) ([]domain.DescriptionTranslation, utils.Pagination, error) {
	return s.translationRepository.GetTranslationsByUser(ctx, userID, pageNumber, pageSize)
}

func (s *DescriptionTranslationService) GetAnimeTranslation(ctx context.Context, animeID int) (*domain.DescriptionTranslation, error) {
	t, err := s.translationRepository.GetTranslationByAnime(ctx, animeID)
	if err != nil {
		return nil, domain_errors.TranslationNotFoundError{AnimeID: strconv.Itoa(animeID)}
	}
	return t, nil
}

func (s *DescriptionTranslationService) GetPendingTranslations(ctx context.Context, pageNumber, pageSize int) ([]domain.DescriptionTranslation, utils.Pagination, error) {
	return s.translationRepository.GetPendingTranslations(ctx, pageNumber, pageSize)
}

func (s *DescriptionTranslationService) AcceptTranslation(ctx context.Context, id int, moderatorID int) error {

	t, err := s.translationRepository.GetTranslationByID(ctx, id)
	if err != nil || t == nil {
		return domain_errors.TranslationNotFoundError{}
	}

	if !t.IsPending() {
		return domain_errors.TranslationNotPendingError{}
	}

	t.Accept(moderatorID)
	return s.translationRepository.UpdateTranslationStatus(ctx, id, t.TranslationStatus, t.AcceptedBy)
}

func (s *DescriptionTranslationService) RejectTranslation(ctx context.Context, id int, moderatorID int) error {

	t, err := s.translationRepository.GetTranslationByID(ctx, id)
	if err != nil || t == nil {
		return domain_errors.TranslationNotFoundError{}
	}

	if !t.IsPending() {
		return domain_errors.TranslationNotPendingError{}
	}

	return s.translationRepository.DeleteTranslation(ctx, id)
}
