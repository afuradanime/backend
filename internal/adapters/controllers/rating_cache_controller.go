package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/go-fuego/fuego"
)

type RatingCacheController struct {
	ratingCacheService interfaces.RatingCacheService
}

func NewRatingCacheController(ratingCacheService interfaces.RatingCacheService) *RatingCacheController {
	return &RatingCacheController{ratingCacheService: ratingCacheService}
}

func (c *RatingCacheController) GetRatingCache(ctx fuego.ContextNoBody) (*domain.RatingCache, error) {
	animeIdStr := ctx.PathParam("animeId")
	animeId, err := strconv.Atoi(animeIdStr)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID: " + err.Error()}
	}

	cache, err := c.ratingCacheService.GetRatingCache(animeId)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to fetch rating cache: " + err.Error()}
	}

	return cache, nil
}
