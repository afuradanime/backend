package controllers

import (
	"errors"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/dtos"
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/go-fuego/fuego"
)

type AnimeListController struct {
	animeListService      interfaces.AnimeListService
	recommendationService interfaces.RecommendationService
}

func NewAnimeListController(s interfaces.AnimeListService, recommendationService interfaces.RecommendationService) *AnimeListController {
	return &AnimeListController{animeListService: s, recommendationService: recommendationService}
}

type UserAnimeListResponse struct {
	Data *dtos.UserAnimeListDTO `json:"data"`
}

func (c *AnimeListController) GetUserList(ctx fuego.ContextNoBody) (UserAnimeListResponse, error) {
	userID, err := strconv.Atoi(ctx.PathParam("userId"))
	if err != nil {
		return UserAnimeListResponse{}, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	var statusFilter *value.AnimeListItemStatus
	statusQuery := ctx.QueryParam("status")
	if statusQuery != "" {
		statusQueryInt, err := strconv.Atoi(statusQuery)
		if err != nil {
			return UserAnimeListResponse{}, fuego.BadRequestError{Detail: "Invalid status filter"}
		}

		st := value.AnimeListItemStatus(statusQueryInt)
		statusFilter = &st
	}

	list, err := c.animeListService.FetchUserList(ctx.Context(), userID, statusFilter)
	if err != nil {
		return UserAnimeListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return UserAnimeListResponse{Data: list}, nil
}

type AddAnimeBody struct {
	Status value.AnimeListItemStatus `json:"status"`
}

type AddAnimeResponse struct {
	Data *dtos.UserListItemDTO `json:"data"`
}

func (c *AnimeListController) AddAnime(ctx fuego.ContextWithBody[AddAnimeBody]) (AddAnimeResponse, error) {

	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return AddAnimeResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.ParseUint(ctx.PathParam("animeId"), 10, 32)
	if err != nil {
		return AddAnimeResponse{}, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return AddAnimeResponse{}, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	status := value.AnimeListItemStatus(body.Status)

	dto, err := c.animeListService.AddAnimeToList(ctx.Context(), userID, uint32(animeID), status)
	if err != nil {
		var animeAlreadyInListErr *domain_errors.AnimeAlreadyInListError
		if errors.As(err, &animeAlreadyInListErr) {
			return AddAnimeResponse{}, fuego.BadRequestError{Detail: err.Error()}
		}

		return AddAnimeResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	// Dismiss recommendation if present
	// We should not do this here, but there was a nasty circular import in the service
	if has, err := c.recommendationService.HasBeenRecommended(ctx.Context(), userID, int(animeID)); err == nil && has {
		c.recommendationService.DismissRecommendation(ctx.Context(), userID, int(animeID))
	}

	return AddAnimeResponse{Data: dto}, nil
}

type UpdateProgressBody struct {
	EpisodesWatched uint32 `json:"episodesWatched"`
}

func (c *AnimeListController) UpdateProgress(ctx fuego.ContextWithBody[UpdateProgressBody]) (any, error) {

	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.ParseUint(ctx.PathParam("animeId"), 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	err = c.animeListService.UpdateProgress(ctx.Context(), userID, uint32(animeID), body.EpisodesWatched)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: err.Error()}
	}

	return nil, nil
}

type UpdateStatusBody struct {
	Status value.AnimeListItemStatus `json:"status"`
}

func (c *AnimeListController) UpdateStatus(ctx fuego.ContextWithBody[UpdateStatusBody]) (any, error) {

	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.ParseUint(ctx.PathParam("animeId"), 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	err = c.animeListService.UpdateStatus(ctx.Context(), userID, uint32(animeID), body.Status)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: err.Error()}
	}

	return nil, nil
}

type UpdateNotesBody struct {
	Notes string `json:"notes"`
}

func (c *AnimeListController) UpdateNotes(ctx fuego.ContextWithBody[UpdateNotesBody]) (any, error) {

	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.ParseUint(ctx.PathParam("animeId"), 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	err = c.animeListService.UpdateNotes(ctx.Context(), userID, uint32(animeID), body.Notes)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: err.Error()}
	}

	return nil, nil
}

type UpdateRatingBody struct {
	Story      uint8 `json:"story"`
	Visuals    uint8 `json:"visuals"`
	Soundtrack uint8 `json:"soundtrack"`
}

func (c *AnimeListController) UpdateRating(ctx fuego.ContextWithBody[UpdateRatingBody]) (any, error) {

	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.ParseUint(ctx.PathParam("animeId"), 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	err = c.animeListService.UpdateRating(
		ctx.Context(),
		userID,
		uint32(animeID),
		body.Story,
		body.Visuals,
		body.Soundtrack,
	)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: err.Error()}
	}

	return nil, nil
}

func (c *AnimeListController) RemoveAnimeFromList(ctx fuego.ContextNoBody) (any, error) {

	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.ParseUint(ctx.PathParam("animeId"), 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	err = c.animeListService.RemoveAnimeFromList(ctx.Context(), userID, uint32(animeID))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: err.Error()}
	}

	return nil, nil
}
