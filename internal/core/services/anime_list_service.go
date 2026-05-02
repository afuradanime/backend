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
	listRepo           	interfaces.AnimeListRepository
	animeRepo          	interfaces.AnimeRepository
	ratingCacheService 	interfaces.RatingCacheService
	userRepo 			interfaces.UserRepository
	mapper             	*mappers.AnimeListMapper
}

func NewAnimeListService(listRepo interfaces.AnimeListRepository, animeRepo interfaces.AnimeRepository, 
	ratingCacheService interfaces.RatingCacheService, userRepo interfaces.UserRepository) *AnimeListService {
	return &AnimeListService{
		listRepo:           listRepo,
		animeRepo:          animeRepo,
		ratingCacheService: ratingCacheService,
		userRepo: 			userRepo,
		mapper:             mappers.NewAnimeListMapper(),
	}
}

func (s *AnimeListService) AddAnimeToList(ctx context.Context, userID int, animeID uint32, status value.AnimeListItemStatus) (*dtos.UserListItemDTO, error) {
	list, err := s.getOrCreateUserList(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Checkpoint 1 - User already has this anime in their list?
	if _, exists := list.GetListItem(animeID); exists {
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

	list.AddListItem(*newItem)

	if err := s.listRepo.SaveUserList(ctx, list); err != nil {
		return nil, err
	}

	return s.mapper.ToDto(newItem, anime), nil
}

func (s *AnimeListService) RemoveAnimeFromList(ctx context.Context, userID int, animeID uint32) error {
	list, err := s.getUserListWithItem(ctx, userID, animeID)
	if err != nil {
		return err
	}

	list.RemoveListItem(animeID)

	return s.listRepo.SaveUserList(ctx, list)
}

func (s *AnimeListService) UpdateProgress(ctx context.Context, userID int, animeID uint32, episodesWatched uint32) error {
	list, item, err := s.getListAndItem(ctx, userID, animeID)
	if err != nil {
		return err
	}

	anime, err := s.animeRepo.FetchAnimeByID(animeID)
	if err != nil || anime == nil {
		return errors.New("failed to fetch anime metadata for validation")
	}

	if err := item.UpdateProgress(episodesWatched, anime.Episodes); err != nil {
		return err
	}

	return s.listRepo.SaveUserList(ctx, list)
}

func (s *AnimeListService) UpdateStatus(ctx context.Context, userID int, animeID uint32, newStatus value.AnimeListItemStatus) error {
	list, item, err := s.getListAndItem(ctx, userID, animeID)
	if err != nil {
		return err
	}

	if item.Status == newStatus {
		return nil
	}

	item.UpdateStatus(newStatus)

	if newStatus == value.AnimeListItemStatusCompleted {
		anime, _ := s.animeRepo.FetchAnimeByID(animeID)
		if anime != nil && anime.Episodes > 0 {
			_ = item.UpdateProgress(anime.Episodes, anime.Episodes)
		}
	}

	return s.listRepo.SaveUserList(ctx, list)
}

func (s *AnimeListService) UpdateNotes(ctx context.Context, userID int, animeID uint32, notes string) error {
	list, item, err := s.getListAndItem(ctx, userID, animeID)
	if err != nil {
		return err
	}

	if err := item.UpdateNotes(notes); err != nil {
		return err
	}

	return s.listRepo.SaveUserList(ctx, list)
}

func (s *AnimeListService) UpdateRating(ctx context.Context, userID int, animeID uint32, story, visuals, soundtrack uint8) error {
	list, item, err := s.getListAndItem(ctx, userID, animeID)
	if err != nil {
		return err
	}

	hadRating := item.Rating != nil
	var oldRating domain.Rating
	if hadRating {
		oldRating = *domain.Uint16ToRating(*item.Rating)
	}

	if err := item.AddRating(story, visuals, soundtrack); err != nil {
		return err
	}

	if !hadRating {
		err = s.ratingCacheService.InsertOrUpdateRating(ctx, userID, int(animeID), story, visuals, soundtrack)
	} else {
		err = s.ratingCacheService.UpdateExistingRating(
			ctx,
			userID, int(animeID),
			oldRating.Story, oldRating.Visuals, oldRating.Soundtrack,
			story, visuals, soundtrack,
		)
	}
	if err != nil {
		log.Printf("Failed to cache rating for user %d and anime %d: %v", userID, animeID, err)
	}

	return s.listRepo.SaveUserList(ctx, list)
}

func (s *AnimeListService) RemoveRating(ctx context.Context, userID int, animeID uint32) error {
	list, item, err := s.getListAndItem(ctx, userID, animeID)
	if err != nil {
		return err
	}

	// Remove from cache, can fail silently if cache is unavailable, but it better not
	if item.Rating != nil {
		readeableRating := domain.Uint16ToRating(*item.Rating)
		s.ratingCacheService.RemoveRating(
			ctx,
			userID,
			int(animeID),
			readeableRating.Story,
			readeableRating.Visuals,
			readeableRating.Soundtrack,
		)
	} else {
		log.Printf("Unexpected state: trying to remove rating for user %d and anime %d but item has no rating", userID, animeID)
	}

	item.RemoveRating()

	return s.listRepo.SaveUserList(ctx, list)
}

func (s *AnimeListService) FetchUserListItem(ctx context.Context, userID int, animeID uint32) (*dtos.UserListItemDTO, error) {
	_, item, err := s.getListAndItem(ctx, userID, animeID)
	if err != nil {
		return nil, err
	}

	anime, err := s.animeRepo.FetchAnimeByID(animeID)
	if err != nil || anime == nil {
		return nil, errors.New("failed to fetch anime metadata for item")
	}

	return s.mapper.ToDto(item, anime), nil
}

func (s *AnimeListService) FetchUserList(ctx context.Context, userID int, viewerID *int, status *value.AnimeListItemStatus) (*dtos.UserAnimeListDTO, error) {
	
	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return nil, domain_errors.UserNotFoundError{}
	}

	// Check if list is private
	log.Printf("PrivateAnimeList: %v", user.PrivateAnimeList)
	log.Printf("viewerID: %v", viewerID)
	if viewerID != nil {
		log.Printf("viewerID value: %d, userID: %d, match: %v", *viewerID, userID, *viewerID == userID)
	}
	if user.PrivateAnimeList {
		isOwner := viewerID != nil && *viewerID == userID
		if !isOwner {
			return nil, &domain_errors.PrivateListError{}
		}
	}
	
	list, err := s.listRepo.FetchUserList(ctx, userID)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return &dtos.UserAnimeListDTO{UserListItems: []*dtos.UserListItemDTO{}}, nil
	}

	items := list.UserListItems

	// Filter by status in the service layer
	if status != nil {
		filtered := items[:0]
		for _, item := range items {
			if item.Status == *status {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	result := make([]*dtos.UserListItemDTO, 0, len(items))
	for _, item := range items {
		anime, err := s.animeRepo.FetchAnimeByID(uint32(item.AnimeID))
		if err != nil || anime == nil {
			// log.Printf("Anime %d was not found for user %d (?)", item.AnimeID, userID)
			continue
		}
		item := item // capture loop var for pointer safety
		result = append(result, s.mapper.ToDto(&item, anime))
	}

	return &dtos.UserAnimeListDTO{
		UserID:        list.UserID,
		UserListItems: result,
	}, nil
}

func (s *AnimeListService) IsInAnimeList(ctx context.Context, receiverID int, animeID int) (bool, error) {

	list, err := s.listRepo.FetchUserList(ctx, receiverID)
	if err != nil {
		return false, err
	}

	if list == nil {
		return false, nil
	}

	_, exists := list.GetListItem(uint32(animeID))
	return exists, nil
}

func (s *AnimeListService) getOrCreateUserList(ctx context.Context, userID int) (*domain.UserAnimeList, error) {
	list, err := s.listRepo.FetchUserList(ctx, userID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = domain.NewPersonalAnimeList(userID)
	}
	return list, nil
}

// fetch the list and validates the anime is present in it.
func (s *AnimeListService) getUserListWithItem(ctx context.Context, userID int, animeID uint32) (*domain.UserAnimeList, error) {
	list, err := s.listRepo.FetchUserList(ctx, userID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, &domain_errors.AnimeNotInListError{
			UserID:  strconv.Itoa(userID),
			AnimeID: strconv.Itoa(int(animeID)),
		}
	}
	if _, exists := list.GetListItem(animeID); !exists {
		return nil, &domain_errors.AnimeNotInListError{
			UserID:  strconv.Itoa(userID),
			AnimeID: strconv.Itoa(int(animeID)),
		}
	}
	return list, nil
}

// fetch the list and returns both the list and a pointer to the specific item.
// The pointer is into the list's slice, so mutations to the item are reflected when saving the list.
func (s *AnimeListService) getListAndItem(ctx context.Context, userID int, animeID uint32) (*domain.UserAnimeList, *domain.UserListItem, error) {
	list, err := s.getUserListWithItem(ctx, userID, animeID)
	if err != nil {
		return nil, nil, err
	}
	item, _ := list.GetListItem(animeID)
	return list, item, nil
}
