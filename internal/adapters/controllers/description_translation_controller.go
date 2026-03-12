package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type DescriptionTranslationController struct {
	translationService interfaces.DescriptionTranslationService
}

func NewDescriptionTranslationController(translationService interfaces.DescriptionTranslationService) *DescriptionTranslationController {
	return &DescriptionTranslationController{
		translationService: translationService,
	}
}

type SubmitTranslationBody struct {
	TranslatedDescription string `json:"TranslatedDescription"`
}

func (c *DescriptionTranslationController) SubmitTranslation(ctx fuego.ContextWithBody[SubmitTranslationBody]) (any, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.Atoi(ctx.PathParam("animeID"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	body, err := ctx.Body()
	if err != nil || body.TranslatedDescription == "" {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	if err := c.translationService.SubmitTranslation(ctx.Context(), animeID, body.TranslatedDescription, userID); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

type AnimeTranslationResponse struct {
	Translation *domain.DescriptionTranslation `json:"translation"`
	Translator  *domain.User                   `json:"translator"`
	Accepter    *domain.User                   `json:"accepter"`
}

func (c *DescriptionTranslationController) GetAnimeTranslation(ctx fuego.ContextNoBody) (AnimeTranslationResponse, error) {
	animeID, err := strconv.Atoi(ctx.PathParam("animeID"))
	if err != nil {
		return AnimeTranslationResponse{}, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	translation, translator, accepter, err := c.translationService.GetAnimeTranslation(ctx.Context(), animeID)
	if err != nil {
		return AnimeTranslationResponse{}, fuego.NotFoundError{Detail: err.Error()}
	}

	return AnimeTranslationResponse{
		Translation: translation,
		Translator:  translator,
		Accepter:    accepter,
	}, nil
}

type PendingTranslationsResponse struct {
	Data       []repositories.PendingTranslationResult `json:"data"`
	Pagination utils.Pagination                        `json:"pagination"`
}

func (c *DescriptionTranslationController) GetPendingTranslations(ctx fuego.ContextNoBody) (PendingTranslationsResponse, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return PendingTranslationsResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	results, pagination, err := c.translationService.GetPendingTranslations(ctx.Context(), pageNumber, pageSize)
	if err != nil {
		return PendingTranslationsResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return PendingTranslationsResponse{
		Data:       results,
		Pagination: pagination,
	}, nil
}

type UserTranslationsResponse struct {
	Data       []domain.DescriptionTranslation `json:"data"`
	Pagination utils.Pagination                `json:"pagination"`
}

func (c *DescriptionTranslationController) GetUserTranslations(ctx fuego.ContextNoBody) (UserTranslationsResponse, error) {
	userID, err := strconv.Atoi(ctx.PathParam("userID"))
	if err != nil {
		return UserTranslationsResponse{}, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)
	translations, pagination, err := c.translationService.GetMyTranslations(ctx.Context(), userID, pageNumber, pageSize)
	if err != nil {
		return UserTranslationsResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return UserTranslationsResponse{
		Data:       translations,
		Pagination: pagination,
	}, nil
}

func (c *DescriptionTranslationController) AcceptTranslation(ctx fuego.ContextNoBody) (any, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	mod, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid translation ID"}
	}

	if err := c.translationService.AcceptTranslation(ctx.Context(), id, mod); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

func (c *DescriptionTranslationController) RejectTranslation(ctx fuego.ContextNoBody) (any, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	mod, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid translation ID"}
	}

	if err := c.translationService.RejectTranslation(ctx.Context(), id, mod); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}
