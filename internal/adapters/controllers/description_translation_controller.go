package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type DescriptionTranslationController struct {
	translationService interfaces.DescriptionTranslationService
}

func NewDescriptionTranslationController(translationService interfaces.DescriptionTranslationService) *DescriptionTranslationController {
	return &DescriptionTranslationController{
		translationService: translationService,
	}
}

func (c *DescriptionTranslationController) SubmitTranslation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	animeStr := chi.URLParam(r, "animeID")
	animeID, err := strconv.Atoi(animeStr)
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	var body struct {
		TranslatedDescription string `json:"TranslatedDescription"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.TranslatedDescription == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.translationService.SubmitTranslation(r.Context(), animeID, body.TranslatedDescription, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *DescriptionTranslationController) GetAnimeTranslation(w http.ResponseWriter, r *http.Request) {
	animeStr := chi.URLParam(r, "animeID")
	animeID, err := strconv.Atoi(animeStr)
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	translation, err := c.translationService.GetAnimeTranslation(r.Context(), animeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(translation)
}

func (c *DescriptionTranslationController) GetMyTranslations(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	translations, pagination, err := c.translationService.GetMyTranslations(r.Context(), userID, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       translations,
		"pagination": pagination,
	})
}

func (c *DescriptionTranslationController) GetPendingTranslations(w http.ResponseWriter, r *http.Request) {

	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	translations, pagination, err := c.translationService.GetPendingTranslations(r.Context(), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       translations,
		"pagination": pagination,
	})
}

func (c *DescriptionTranslationController) AcceptTranslation(w http.ResponseWriter, r *http.Request) {
	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mod, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid translation ID", http.StatusBadRequest)
		return
	}

	if err := c.translationService.AcceptTranslation(r.Context(), id, mod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *DescriptionTranslationController) RejectTranslation(w http.ResponseWriter, r *http.Request) {
	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mod, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid translation ID", http.StatusBadRequest)
		return
	}

	if err := c.translationService.RejectTranslation(r.Context(), id, mod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
