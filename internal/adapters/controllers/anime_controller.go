package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/core/interfaces"
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
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Get pagination parameters with defaults
	pageNumber := 0
	pageSize := 50 // Im keeping page sized fixed, maybe move it to the config or let the user set it idk

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	animes, err := ac.animeService.FetchAnimeFromQuery(query, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}

func (ac *AnimeController) GetAnimeThisSeason(w http.ResponseWriter, r *http.Request) {
	animes, err := ac.animeService.FetchAnimeThisSeason()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}
