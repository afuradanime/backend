package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/filters"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type AnimeController struct {
	animeService interfaces.AnimeService
}

func NewAnimeController(s interfaces.AnimeService) *AnimeController {
	return &AnimeController{animeService: s}
}

func parseAnimeFilters(ctx fuego.ContextNoBody) filters.AnimeFilter {
	var f filters.AnimeFilter

	if name := ctx.QueryParam("q"); name != "" {
		f.Name = &name
	}
	if typeStr := ctx.QueryParam("type"); typeStr != "" {
		if t, err := strconv.ParseUint(typeStr, 10, 32); err == nil {
			t32 := uint32(t)
			f.Type = &t32
		}
	}
	if statusStr := ctx.QueryParam("status"); statusStr != "" {
		if s, err := strconv.ParseUint(statusStr, 10, 32); err == nil {
			s32 := uint32(s)
			f.Status = &s32
		}
	}
	if startStr := ctx.QueryParam("start_date"); startStr != "" {
		if t, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			f.StartDate = &t
		}
	}
	if endStr := ctx.QueryParam("end_date"); endStr != "" {
		if t, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			f.EndDate = &t
		}
	}
	if minStr := ctx.QueryParam("min_episodes"); minStr != "" {
		if m, err := strconv.ParseUint(minStr, 10, 32); err == nil {
			m32 := uint32(m)
			f.MinEpisodes = &m32
		}
	}
	if maxStr := ctx.QueryParam("max_episodes"); maxStr != "" {
		if m, err := strconv.ParseUint(maxStr, 10, 32); err == nil {
			m32 := uint32(m)
			f.MaxEpisodes = &m32
		}
	}

	return f
}

func (ac *AnimeController) GetAnimeByID(ctx fuego.ContextNoBody) (*domain.Anime, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid anime ID"}
	}

	anime, err := ac.animeService.FetchAnimeByID(uint32(id))
	if err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return anime, nil
}

func (ac *AnimeController) GetRandomAnime(ctx fuego.ContextNoBody) (*domain.Anime, error) {
	anime, err := ac.animeService.FetchRandomAnime()
	if err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}
	return anime, nil
}

type AnimeListResponse struct {
	Animes     []*domain.Anime  `json:"animes"`
	Pagination utils.Pagination `json:"pagination"`
}

func (ac *AnimeController) SearchAnime(ctx fuego.ContextNoBody) (AnimeListResponse, error) {
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)
	f := parseAnimeFilters(ctx)

	animes, pagination, err := ac.animeService.FetchAnimeFromQuery(f, pageNumber, pageSize)
	if err != nil {
		return AnimeListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return AnimeListResponse{Animes: animes, Pagination: pagination}, nil
}

func (ac *AnimeController) GetAnimeThisSeason(ctx fuego.ContextNoBody) (AnimeListResponse, error) {
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)
	f := parseAnimeFilters(ctx)

	animes, pagination, err := ac.animeService.FetchAnimeThisSeason(f, pageNumber, pageSize)
	if err != nil {
		return AnimeListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return AnimeListResponse{Animes: animes, Pagination: pagination}, nil
}

type StudioAnimeResponse struct {
	Studio     *value.Studio    `json:"studio"`
	Animes     []*domain.Anime  `json:"animes"`
	Pagination utils.Pagination `json:"pagination"`
}

func (ac *AnimeController) GetAnimeByStudioID(ctx fuego.ContextNoBody) (StudioAnimeResponse, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return StudioAnimeResponse{}, fuego.BadRequestError{Detail: "Invalid studio ID"}
	}
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)
	f := parseAnimeFilters(ctx)

	studio, animes, pagination, err := ac.animeService.FetchStudioByID(f, uint32(id), pageNumber, pageSize)
	if err != nil {
		return StudioAnimeResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return StudioAnimeResponse{Studio: studio, Animes: animes, Pagination: pagination}, nil
}

type ProducerAnimeResponse struct {
	Producer   *value.Producer  `json:"producer"`
	Animes     []*domain.Anime  `json:"animes"`
	Pagination utils.Pagination `json:"pagination"`
}

func (ac *AnimeController) GetAnimeByProducerID(ctx fuego.ContextNoBody) (ProducerAnimeResponse, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return ProducerAnimeResponse{}, fuego.BadRequestError{Detail: "Invalid producer ID"}
	}
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)
	f := parseAnimeFilters(ctx)

	producer, animes, pagination, err := ac.animeService.FetchProducerByID(f, uint32(id), pageNumber, pageSize)
	if err != nil {
		return ProducerAnimeResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return ProducerAnimeResponse{Producer: producer, Animes: animes, Pagination: pagination}, nil
}

type LicensorAnimeResponse struct {
	Licensor   *value.Licensor  `json:"licensor"`
	Animes     []*domain.Anime  `json:"animes"`
	Pagination utils.Pagination `json:"pagination"`
}

func (ac *AnimeController) GetAnimeByLicensorID(ctx fuego.ContextNoBody) (LicensorAnimeResponse, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return LicensorAnimeResponse{}, fuego.BadRequestError{Detail: "Invalid licensor ID"}
	}
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)
	f := parseAnimeFilters(ctx)

	licensor, animes, pagination, err := ac.animeService.FetchLicensorByID(f, uint32(id), pageNumber, pageSize)
	if err != nil {
		return LicensorAnimeResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return LicensorAnimeResponse{Licensor: licensor, Animes: animes, Pagination: pagination}, nil
}

func (ac *AnimeController) GetAnimeByTagID(ctx fuego.ContextNoBody) (AnimeListResponse, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return AnimeListResponse{}, fuego.BadRequestError{Detail: "Invalid tag ID"}
	}
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)
	f := parseAnimeFilters(ctx)

	animes, pagination, err := ac.animeService.FetchAnimeFromTag(uint32(id), f, pageNumber, pageSize)
	if err != nil {
		return AnimeListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return AnimeListResponse{Animes: animes, Pagination: pagination}, nil
}
