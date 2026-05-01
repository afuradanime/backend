package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type AnimeWithRating struct {
	Anime  *domain.Anime       `json:"anime"`
	Rating *domain.RatingCache `json:"rating"`
}

type PaginatedAnimeWithRating struct {
	Data       []AnimeWithRating `json:"data"`
	Pagination utils.Pagination  `json:"pagination"`
}

type RatingCacheController struct {
	ratingCacheService interfaces.RatingCacheService
	animeService       interfaces.AnimeService
}

func NewRatingCacheController(ratingCacheService interfaces.RatingCacheService, animeService interfaces.AnimeService) *RatingCacheController {
	return &RatingCacheController{
		ratingCacheService: ratingCacheService,
		animeService:       animeService,
	}
}

func (c *RatingCacheController) GetRatingCache(ctx fuego.ContextNoBody) (*domain.RatingCache, error) {
	animeIdStr := ctx.PathParam("animeId")
	animeId, err := strconv.Atoi(animeIdStr)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID: " + err.Error()}
	}

	cache, err := c.ratingCacheService.GetRatingCache(ctx, animeId)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to fetch rating cache: " + err.Error()}
	}

	return cache, nil
}

type PaginatedRatingCache struct {
	Data       []*domain.RatingCache `json:"data"`
	Pagination utils.Pagination      `json:"pagination"`
}

func (c *RatingCacheController) GetTopAnime(ctx fuego.ContextNoBody) (*PaginatedAnimeWithRating, error) {
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 10)

	caches, pagination, err := c.ratingCacheService.GetTopAnime(ctx.Context(), pageNumber, pageSize)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to fetch top anime: " + err.Error()}
	}
	return c.enrichWithAnime(caches, pagination)
}

func (c *RatingCacheController) GetPopularAnime(ctx fuego.ContextNoBody) (*PaginatedAnimeWithRating, error) {
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 10)

	caches, pagination, err := c.ratingCacheService.GetPopularAnime(ctx.Context(), pageNumber, pageSize)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to fetch popular anime: " + err.Error()}
	}
	return c.enrichWithAnime(caches, pagination)
}

func (c *RatingCacheController) enrichWithAnime(caches []*domain.RatingCache, pagination utils.Pagination) (*PaginatedAnimeWithRating, error) {
	data := make([]AnimeWithRating, 0, len(caches))
	for _, cache := range caches {
		anime, err := c.animeService.FetchAnimeByID(uint32(cache.AnimeID))
		if err != nil {
			continue // skip if anime not found
		}
		data = append(data, AnimeWithRating{Anime: anime, Rating: cache})
	}
	return &PaginatedAnimeWithRating{Data: data, Pagination: pagination}, nil
}
