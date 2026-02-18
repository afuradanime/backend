package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type AnimeController struct {
	animeService interfaces.AnimeService
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
	if err != nil { // TODO: Proper error handling here, with different status codes!
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

func (ac *AnimeController) SearchAnime(w http.ResponseWriter, r *http.Request) {
	// Get query parameter
	query := r.URL.Query().Get("q")
	// if query == "" {
	// 	http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
	// 	return
	// }

	// Get pagination parameters with defaults
	pageNumber := 1
	pageSize := 50

	if pageStr := r.URL.Query().Get("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	if pageSize > 50 {
		pageSize = 50
	}

	animes, pagination, err := ac.animeService.FetchAnimeFromQuery(query, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{
		Animes:     animes,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ac *AnimeController) GetAnimeThisSeason(w http.ResponseWriter, r *http.Request) {

	pageNumber := 1
	pageSize := 50

	if pageStr := r.URL.Query().Get("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	if pageSize > 50 {
		pageSize = 50
	}

	animes, pagination, err := ac.animeService.FetchAnimeThisSeason(pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{
		Animes:     animes,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ac *AnimeController) GetAnimeByStudioID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid studio ID", http.StatusBadRequest)
		return
	}

	pageNumber := 1
	pageSize := 50

	if pageStr := r.URL.Query().Get("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	if pageSize > 50 {
		pageSize = 50
	}

	studio, animes, pagination, err := ac.animeService.FetchStudioByID(uint32(id), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Studio     *value.Studio    `json:"studio"`
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{
		Studio:     studio,
		Animes:     animes,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ac *AnimeController) GetAnimeByProducerID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid producer ID", http.StatusBadRequest)
		return
	}

	pageNumber := 1
	pageSize := 50

	if pageStr := r.URL.Query().Get("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	if pageSize > 50 {
		pageSize = 50
	}

	producer, animes, pagination, err := ac.animeService.FetchProducerByID(uint32(id), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Producer   *value.Producer  `json:"producer"`
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{
		Producer:   producer,
		Animes:     animes,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ac *AnimeController) GetAnimeByLicensorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid licensor ID", http.StatusBadRequest)
		return
	}

	pageNumber := 1
	pageSize := 50

	if pageStr := r.URL.Query().Get("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	if pageSize > 50 {
		pageSize = 50
	}

	licensor, animes, pagination, err := ac.animeService.FetchLicensorByID(uint32(id), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Licensor   *value.Licensor  `json:"licensor"`
		Animes     []*domain.Anime  `json:"animes"`
		Pagination utils.Pagination `json:"pagination"`
	}{
		Licensor:   licensor,
		Animes:     animes,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
