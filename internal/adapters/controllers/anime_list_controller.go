package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/go-chi/chi/v5"
)

type AnimeListController struct {
	animeListService interfaces.AnimeListService
}

func NewAnimeListController(s interfaces.AnimeListService) *AnimeListController {
	return &AnimeListController{animeListService: s}
}

// GET /animelist/{userId}?status={status}
func (c *AnimeListController) GetUserList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var statusFilter *value.AnimeListItemStatus
	statusQuery := r.URL.Query().Get("status")
	if statusQuery != "" {
		statusQueryInt, err := strconv.Atoi(statusQuery)
		if err != nil {
			http.Error(w, "Invalid status filter", http.StatusBadRequest)
			return
		}
		st := value.AnimeListItemStatus(statusQueryInt)
		statusFilter = &st
	}

	list, err := c.animeListService.FetchUserList(r.Context(), userID, statusFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// POST /animelist/{userId}/{animeId}
func (c *AnimeListController) AddAnime(w http.ResponseWriter, r *http.Request) {
	userID, animeID, err := parseUserAndAnimeIDs(r)
	if err != nil {
		http.Error(w, "Invalid IDs in URL", http.StatusBadRequest)
		return
	}

	// Check if the user is trying to add to their own list
	allowed := allowedToModifyList(r, userID)
	if !allowed {
		// Another case of a user trying to modify something that doesn't belong to them...
		// We could report them internally again for abusing our API, but for now let's just return Unauthorized
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Status value.AnimeListItemStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := value.AnimeListItemStatus(body.Status)

	dto, err := c.animeListService.AddAnimeToList(r.Context(), userID, animeID, status)
	if err != nil {
		var animeAlreadyInListErr *domain_errors.AnimeAlreadyInListError
		if errors.As(err, &animeAlreadyInListErr) {
			http.Error(w, "Anime already in list", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 = Created
	json.NewEncoder(w).Encode(dto)
}

// PATCH /animelist/{userId}/progress/{animeId}
func (c *AnimeListController) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	userID, animeID, err := parseUserAndAnimeIDs(r)
	if err != nil {
		http.Error(w, "Invalid IDs in URL", http.StatusBadRequest)
		return
	}

	allowed := allowedToModifyList(r, userID)
	if !allowed {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		EpisodesWatched uint32 `json:"episodesWatched"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.animeListService.UpdateProgress(r.Context(), userID, animeID, body.EpisodesWatched)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /animelist/{userId}/status/{animeId}
func (c *AnimeListController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, animeID, err := parseUserAndAnimeIDs(r)
	if err != nil {
		http.Error(w, "Invalid IDs in URL", http.StatusBadRequest)
		return
	}

	allowed := allowedToModifyList(r, userID)
	if !allowed {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Status value.AnimeListItemStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := value.AnimeListItemStatus(body.Status)
	err = c.animeListService.UpdateStatus(r.Context(), userID, animeID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /animelist/{userId}/notes/{animeId}
func (c *AnimeListController) UpdateNotes(w http.ResponseWriter, r *http.Request) {
	userID, animeID, err := parseUserAndAnimeIDs(r)
	if err != nil {
		http.Error(w, "Invalid IDs in URL", http.StatusBadRequest)
		return
	}

	allowed := allowedToModifyList(r, userID)
	if !allowed {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Notes string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.animeListService.UpdateNotes(r.Context(), userID, animeID, body.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /animelist/{userId}/rating/{animeId}
func (c *AnimeListController) UpdateRating(w http.ResponseWriter, r *http.Request) {
	userID, animeID, err := parseUserAndAnimeIDs(r)
	if err != nil {
		http.Error(w, "Invalid IDs in URL", http.StatusBadRequest)
		return
	}

	allowed := allowedToModifyList(r, userID)
	if !allowed {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Story      uint8 `json:"story"`
		Visuals    uint8 `json:"visuals"`
		Soundtrack uint8 `json:"soundtrack"`
		Enjoyment  uint8 `json:"enjoyment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.animeListService.UpdateRating(r.Context(), userID, animeID, body.Story, body.Visuals, body.Soundtrack, body.Enjoyment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /animelist/{userId}/{animeId}
func (c *AnimeListController) RemoveAnimeFromList(w http.ResponseWriter, r *http.Request) {
	userID, animeID, err := parseUserAndAnimeIDs(r)
	if err != nil {
		http.Error(w, "Invalid IDs in URL", http.StatusBadRequest)
		return
	}

	allowed := allowedToModifyList(r, userID)
	if !allowed {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = c.animeListService.RemoveAnimeFromList(r.Context(), userID, animeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseUserAndAnimeIDs(r *http.Request) (int, uint32, error) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		return 0, 0, err
	}

	animeIDStr := chi.URLParam(r, "animeId")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		return 0, 0, err
	}

	return userID, uint32(animeID), nil
}

func allowedToModifyList(r *http.Request, targetUserID int) bool {
	loggedUserID, ok := middlewares.GetUserIDFromContext(r)
	if !ok || loggedUserID != targetUserID {
		return false
	}
	return true
}
