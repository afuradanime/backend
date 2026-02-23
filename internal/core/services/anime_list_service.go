package services

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/dtos"
	"github.com/afuradanime/backend/internal/adapters/mappers"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type AnimeListService struct {
	listRepo  interfaces.AnimeListRepository
	animeRepo interfaces.AnimeRepository
	mapper    *mappers.AnimeListMapper
}

func NewAnimeListService(listRepo interfaces.AnimeListRepository, animeRepo interfaces.AnimeRepository) *AnimeListService {
	return &AnimeListService{
		listRepo:  listRepo,
		animeRepo: animeRepo,
		mapper:    mappers.NewAnimeListMapper(),
	}
}

func (s *AnimeListService) AddAnimeToList(ctx context.Context, userID int, animeID uint32, status value.AnimeListItemStatus) (*dtos.AnimeListItemDTO, error) {
	// Checkpoint 1 - User already has this anime in their list?
	item, err := getAndValidateAnimeInList(ctx, s.listRepo, userID, animeID)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return nil, &domain_errors.AnimeAlreadyInListError{
			UserID:  strconv.Itoa(userID),
			AnimeID: strconv.Itoa(int(animeID)),
		}
	}

	// Checkpoint 2 - Anime exists ?
	anime, err := s.animeRepo.FetchAnimeByID(animeID)
	if err != nil {
		return nil, err
	}
	if anime == nil {
		return nil, domain_errors.AnimeNotFoundError{
			AnimeID: strconv.Itoa(int(animeID)),
		}
	}

	// All checks passed, we can create the AnimeListItem domain object
	newItem := domain.NewAnimeListItem(userID, animeID, status)

	// User marked as completed, therefore they watched all episodes
	if status == value.AnimeListItemStatusCompleted && anime.Episodes > 0 {
		_ = newItem.UpdateProgress(anime.Episodes, anime.Episodes)
	}

	if err := s.listRepo.AddListItem(ctx, newItem); err != nil {
		return nil, err
	}

	return s.mapper.ToDto(newItem, anime), nil
}

func (s *AnimeListService) UpdateProgress(ctx context.Context, userID int, animeID uint32, episodesWatched uint32) error {
	// Get list item
	item, err := getAndValidateAnimeInList(ctx, s.listRepo, userID, animeID)
	if err != nil {
		return err
	}

	// Get anime
	anime, err := s.animeRepo.FetchAnimeByID(animeID)
	if err != nil || anime == nil {
		return errors.New("failed to fetch anime metadata for validation")
	}

	if err := item.UpdateProgress(episodesWatched, anime.Episodes); err != nil {
		return err
	}

	return s.listRepo.UpdateListItem(ctx, item)
}

func (s *AnimeListService) UpdateStatus(ctx context.Context, userID int, animeID uint32, newStatus value.AnimeListItemStatus) error {
	// Get list item
	item, err := getAndValidateAnimeInList(ctx, s.listRepo, userID, animeID)
	if err != nil {
		return err
	}

	// Same status as before, no update needed
	if item.Status == newStatus {
		return nil
	}

	item.UpdateStatus(newStatus)

	// completed therefore they watched all episodes
	if newStatus == value.AnimeListItemStatusCompleted {
		anime, _ := s.animeRepo.FetchAnimeByID(animeID)
		if anime != nil && anime.Episodes > 0 {
			_ = item.UpdateProgress(anime.Episodes, anime.Episodes)
		}
	}

	return s.listRepo.UpdateListItem(ctx, item)
}

func (s *AnimeListService) UpdateNotes(ctx context.Context, userID int, animeID uint32, notes string) error {
	// Get list item
	item, err := getAndValidateAnimeInList(ctx, s.listRepo, userID, animeID)
	if err != nil {
		return err
	}

	if err := item.UpdateNotes(notes); err != nil {
		return err
	}

	return s.listRepo.UpdateListItem(ctx, item)
}

func (s *AnimeListService) UpdateRating(ctx context.Context, userID int, animeID uint32, story, visuals, soundtrack, enjoyment uint8) error {
	// Get list item
	item, err := getAndValidateAnimeInList(ctx, s.listRepo, userID, animeID)
	if err != nil {
		return err
	}

	if err := item.AddRating(story, visuals, soundtrack, enjoyment); err != nil {
		return err
	}

	return s.listRepo.UpdateListItem(ctx, item)
}

func (s *AnimeListService) RemoveAnimeFromList(ctx context.Context, userID int, animeID uint32) error {
	return s.listRepo.DeleteListItem(ctx, userID, animeID)
}

func (s *AnimeListService) FetchUserListItem(ctx context.Context, userID int, animeID uint32) (*dtos.AnimeListItemDTO, error) {
	item, err := getAndValidateAnimeInList(ctx, s.listRepo, userID, animeID)
	if err != nil {
		return nil, err
	}

	anime, err := s.animeRepo.FetchAnimeByID(animeID)
	if err != nil || anime == nil {
		return nil, errors.New("failed to fetch anime metadata for item")
	}

	return s.mapper.ToDto(item, anime), nil
}

func (s *AnimeListService) FetchUserList(ctx context.Context, userID int, status *value.AnimeListItemStatus) ([]*dtos.AnimeListItemDTO, error) {
	listItems, err := s.listRepo.FetchUserList(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	if len(listItems) == 0 {
		return make([]*dtos.AnimeListItemDTO, 0), nil
	}

	result := make([]*dtos.AnimeListItemDTO, 0, len(listItems))

	for _, item := range listItems {

		anime, err := s.animeRepo.FetchAnimeByID(item.AnimeID)
		if err != nil || anime == nil {
			log.Printf("Anime %d was not found for user %d (?)", item.AnimeID, userID)
			continue
		}

		dto := s.mapper.ToDto(item, anime)
		result = append(result, dto)
	}

	return result, nil
}

func getAndValidateAnimeInList(ctx context.Context, listRepo interfaces.AnimeListRepository, userID int, animeID uint32) (*domain.AnimeListItem, error) {
	item, err := listRepo.FetchItemByUserAndAnime(ctx, userID, animeID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, &domain_errors.AnimeNotInListError{
			UserID:  strconv.Itoa(userID),
			AnimeID: strconv.Itoa(int(animeID)),
		}
	}
	return item, nil
}
