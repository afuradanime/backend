package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(s interfaces.UserService) *UserController {
	return &UserController{userService: s}
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	users, pagination, err := uc.userService.GetUsers(r.Context(), pageNumber, pageSize)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       users,
		"pagination": pagination,
	})
}

func (uc *UserController) SearchByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("q")
	if username == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	users, pagination, err := uc.userService.SearchByUsername(r.Context(), username, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       users,
		"pagination": pagination,
	})
}

func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := uc.userService.GetUserByID(r.Context(), id)
	if err != nil { // TODO: Proper error handling here, with different status codes!
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	id, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateData struct {
		Email                 *string   `json:"Email"`
		Username              *string   `json:"Username"`
		Location              *string   `json:"Location"`
		Pronouns              *string   `json:"Pronouns"`
		Socials               *[]string `json:"Socials"`
		Birthday              *string   `json:"Birthday"`
		AllowsFriendRequests  *bool     `json:"AllowsFriendRequests"`
		AllowsRecommendations *bool     `json:"AllowsRecommendations"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var birthday *time.Time
	if updateData.Birthday != nil {
		t, err := time.Parse("2006-01-02", *updateData.Birthday)
		if err != nil {
			http.Error(w, "Invalid birthday format, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		birthday = &t
	}

	err := uc.userService.UpdatePersonalInfo(
		r.Context(), id,
		updateData.Email,
		updateData.Username,
		updateData.Location,
		updateData.Pronouns,
		updateData.Socials,
		birthday,
		updateData.AllowsFriendRequests,
		updateData.AllowsRecommendations,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uc *UserController) RestrictAccount(w http.ResponseWriter, r *http.Request) {
	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var body struct {
		CanPost      bool `json:"CanPost"`
		CanTranslate bool `json:"CanTranslate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uc.userService.RestrictAccount(r.Context(), targetID, body.CanPost, body.CanTranslate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
