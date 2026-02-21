package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/filters"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type AnimeController struct {
	animeService interfaces.AnimeService
}

func parseAnimeFilters(r *http.Request) filters.AnimeFilter {
	var f filters.AnimeFilter

	if name := r.URL.Query().Get("q"); name != "" {
		f.Name = &name
	}
	if typeStr := r.URL.Query().Get("type"); typeStr != "" {
		if t, err := strconv.ParseUint(typeStr, 10, 32); err == nil {
			t32 := uint32(t)
			f.Type = &t32
		}
	}
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		if s, err := strconv.ParseUint(statusStr, 10, 32); err == nil {
			s32 := uint32(s)
			f.Status = &s32
		}
	}
	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		if t, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			f.StartDate = &t
		}
	}
	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		if t, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			f.EndDate = &t
		}
	}
	if minStr := r.URL.Query().Get("min_episodes"); minStr != "" {
		if m, err := strconv.ParseUint(minStr, 10, 32); err == nil {
			m32 := uint32(m)
			f.MinEpisodes = &m32
		}
	}
	if maxStr := r.URL.Query().Get("max_episodes"); maxStr != "" {
		if m, err := strconv.ParseUint(maxStr, 10, 32); err == nil {
			m32 := uint32(m)
			f.MaxEpisodes = &m32
		}
	}

	return f
}

func NewAnimeController(s interfaces.AnimeService) *AnimeController {
	return &AnimeController{animeService: s}
}

func (ac *AnimeController) GetAnimeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	anime, err := ac.animeService.FetchAnimeByID(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

func (ac *AnimeController) GetRandomAnime(w http.ResponseWriter, r *http.Request) {

	anime, err := ac.animeService.FetchRandomAnime()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

func (ac *AnimeController) SearchAnime(w http.ResponseWriter, r *http.Request) {
	pageNumber, pageSize := utils.GetPaginationParams(r, 50)
	f := parseAnimeFilters(r)

	animes, pagination, err := ac.animeService.FetchAnimeFromQuery(f, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{animes, pagination})
}

func (ac *AnimeController) GetAnimeThisSeason(w http.ResponseWriter, r *http.Request) {
	pageNumber, pageSize := utils.GetPaginationParams(r, 50)
	f := parseAnimeFilters(r)

	animes, pagination, err := ac.animeService.FetchAnimeThisSeason(f, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{animes, pagination})
}

func (ac *AnimeController) GetAnimeByStudioID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid studio ID", http.StatusBadRequest)
		return
	}
	pageNumber, pageSize := utils.GetPaginationParams(r, 50)
	f := parseAnimeFilters(r)

	studio, animes, pagination, err := ac.animeService.FetchStudioByID(f, uint32(id), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Studio     *value.Studio    `json:"studio"`
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{studio, animes, pagination})
}

func (ac *AnimeController) GetAnimeByProducerID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid producer ID", http.StatusBadRequest)
		return
	}
	pageNumber, pageSize := utils.GetPaginationParams(r, 50)
	f := parseAnimeFilters(r)

	producer, animes, pagination, err := ac.animeService.FetchProducerByID(f, uint32(id), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Producer   *value.Producer  `json:"producer"`
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{producer, animes, pagination})
}

func (ac *AnimeController) GetAnimeByLicensorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid licensor ID", http.StatusBadRequest)
		return
	}
	pageNumber, pageSize := utils.GetPaginationParams(r, 50)
	f := parseAnimeFilters(r)

	licensor, animes, pagination, err := ac.animeService.FetchLicensorByID(f, uint32(id), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Licensor   *value.Licensor  `json:"licensor"`
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{licensor, animes, pagination})
}

func (ac *AnimeController) GetAnimeByTagID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}
	pageNumber, pageSize := utils.GetPaginationParams(r, 50)
	f := parseAnimeFilters(r)

	animes, pagination, err := ac.animeService.FetchAnimeFromTag(uint32(id), f, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{animes, pagination})
}
