package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/interfaces"
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

func getUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
	return userID, ok
}

func (c *FriendshipController) SendFriendRequest(w http.ResponseWriter, r *http.Request) {

	initiator, ok := getUserIDFromContext(r)
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
	receiver, ok := getUserIDFromContext(r)
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
	receiver, ok := getUserIDFromContext(r)
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
	initiator, ok := getUserIDFromContext(r)
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

	friends, pagination, err := c.friendshipService.GetFriendList(r.Context(), targetUser, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return both data and pagination metadata
	resp := map[string]interface{}{
		"data":       friends,
		"pagination": pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *FriendshipController) ListPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	requests, pagination, err := c.friendshipService.GetPendingFriendRequests(r.Context(), userID, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return both data and pagination metadata
	resp := map[string]interface{}{
		"data":       requests,
		"pagination": pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *FriendshipController) AreFriends(w http.ResponseWriter, r *http.Request) {
	userA, ok := getUserIDFromContext(r)
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

	areFriends, err := c.friendshipService.AreFriends(r.Context(), userA, userB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"areFriends": areFriends})
}
