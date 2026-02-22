package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type RecommendationController struct {
	service interfaces.RecommendationService
}

func NewRecommendationController(service interfaces.RecommendationService) *RecommendationController {
	return &RecommendationController{service: service}
}

func (c *RecommendationController) Send(w http.ResponseWriter, r *http.Request) {
	initiatorID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverID, err := strconv.Atoi(chi.URLParam(r, "receiverID"))
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	animeID, err := strconv.Atoi(chi.URLParam(r, "animeID"))
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	err = c.service.Send(r.Context(), initiatorID, receiverID, animeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *RecommendationController) GetMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	recs, pagination, err := c.service.GetUserRecommendations(r.Context(), userID, pageNumber, pageSize)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       recs,
		"pagination": pagination,
	})
}

func (c *RecommendationController) Dismiss(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	animeID, err := strconv.Atoi(chi.URLParam(r, "animeID"))
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DismissRecommendation(r.Context(), userID, animeID); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
