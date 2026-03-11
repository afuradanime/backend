package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type RecommendationController struct {
	service interfaces.RecommendationService
}

func NewRecommendationController(service interfaces.RecommendationService) *RecommendationController {
	return &RecommendationController{service: service}
}

func (c *RecommendationController) Send(ctx fuego.ContextNoBody) (any, error) {
	initiatorID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	receiverID, err := strconv.Atoi(ctx.PathParam("receiverID"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid receiver ID"}
	}

	animeID, err := strconv.Atoi(ctx.PathParam("animeID"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	err = c.service.Send(ctx.Context(), initiatorID, receiverID, animeID)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: err.Error()}
	}

	return nil, nil
}

type UserRecommendationsResponse struct {
	Data       []*domain.Recommendation `json:"data"`
	Pagination utils.Pagination         `json:"pagination"`
}

func (c *RecommendationController) GetMine(ctx fuego.ContextNoBody) (UserRecommendationsResponse, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return UserRecommendationsResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	recs, pagination, err := c.service.GetUserRecommendations(ctx.Context(), userID, pageNumber, pageSize)
	if err != nil {
		return UserRecommendationsResponse{}, fuego.InternalServerError{Detail: "Internal server error"}
	}

	return UserRecommendationsResponse{
		Data:       recs,
		Pagination: pagination,
	}, nil
}

func (c *RecommendationController) Dismiss(ctx fuego.ContextNoBody) (any, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	animeID, err := strconv.Atoi(ctx.PathParam("animeID"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	if err := c.service.DismissRecommendation(ctx.Context(), userID, animeID); err != nil {
		return nil, fuego.InternalServerError{Detail: "Internal server error"}
	}

	return nil, nil
}
