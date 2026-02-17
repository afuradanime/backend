package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// helper to get initiator and receiver safely
func getInitiatorAndReceiver(r *http.Request) (int, int, error) {
	initiatorStr := chi.URLParam(r, "initiator")
	receiverStr := chi.URLParam(r, "receiver")

	if initiatorStr == "" || receiverStr == "" {
		return 0, 0, http.ErrMissingFile
	}

	initiator, err := strconv.Atoi(initiatorStr)
	if err != nil {
		return 0, 0, err
	}

	receiver, err := strconv.Atoi(receiverStr)
	if err != nil {
		return 0, 0, err
	}

	return initiator, receiver, nil
}

func (c *FriendshipController) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	initiator, receiver, err := getInitiatorAndReceiver(r)
	if err != nil {
		http.Error(w, "Both initiator and receiver are required", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.SendFriendRequest(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	initiator, receiver, err := getInitiatorAndReceiver(r)
	if err != nil {
		http.Error(w, "Both initiator and receiver are required", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.AcceptFriendRequest(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
	initiator, receiver, err := getInitiatorAndReceiver(r)
	if err != nil {
		http.Error(w, "Both initiator and receiver are required", http.StatusBadRequest)
		return
	}

	if err := c.friendshipService.DeclineFriendRequest(r.Context(), initiator, receiver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FriendshipController) BlockUser(w http.ResponseWriter, r *http.Request) {
	initiator, receiver, err := getInitiatorAndReceiver(r)
	if err != nil {
		http.Error(w, "Both initiator and receiver are required", http.StatusBadRequest)
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

	friends, err := c.friendshipService.GetFriendList(r.Context(), targetUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friends); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *FriendshipController) ListPendingFriendRequests(w http.ResponseWriter, r *http.Request) {

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

	requests, err := c.friendshipService.GetPendingFriendRequests(r.Context(), targetUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(requests); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
