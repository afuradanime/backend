package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type FriendshipController struct {
	friendshipService interfaces.FriendshipService
}

func NewFriendshipController(friendshipService interfaces.FriendshipService) *FriendshipController {
	return &FriendshipController{
		friendshipService: friendshipService,
	}
}

func (c *FriendshipController) SendFriendRequest(w http.ResponseWriter, r *http.Request) {

	initiator, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverStr := chi.URLParam(r, "receiver")
	receiver, err := strconv.Atoi(receiverStr)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.SendFriendRequest(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	receiver, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	initiatorStr := chi.URLParam(r, "initiator")
	initiator, err := strconv.Atoi(initiatorStr)
	if err != nil {
		http.Error(w, "Invalid initiator ID", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.AcceptFriendRequest(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
	receiver, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	initiatorStr := chi.URLParam(r, "initiator")
	initiator, err := strconv.Atoi(initiatorStr)
	if err != nil {
		http.Error(w, "Invalid initiator ID", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.DeclineFriendRequest(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) BlockUser(w http.ResponseWriter, r *http.Request) {
	initiator, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverStr := chi.URLParam(r, "receiver")
	receiver, err := strconv.Atoi(receiverStr)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.BlockUser(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) ListFriends(w http.ResponseWriter, r *http.Request) {
	targetUserStr := chi.URLParam(r, "userID")
	if targetUserStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	targetUser, err := strconv.Atoi(targetUserStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 50)

	friends, pagination, err := c.friendshipService.GetFriendList(r.Context(), targetUser, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       friends,
		"pagination": pagination,
	})
}

func (c *FriendshipController) ListPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 50)

	requests, pagination, err := c.friendshipService.GetPendingFriendRequests(r.Context(), userID, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       requests,
		"pagination": pagination,
	})
}

func (c *FriendshipController) FetchFriendshipStatus(w http.ResponseWriter, r *http.Request) {
	userA, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverStr := chi.URLParam(r, "receiver")
	userB, err := strconv.Atoi(receiverStr)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	friendshipStatus, err := c.friendshipService.FetchFriendshipStatus(r.Context(), userA, userB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if friendshipStatus == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.Friendship{
			Initiator: userA,
			Receiver:  userB,
			Status:    value.FriendshipStatusNone,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*friendshipStatus)
}
