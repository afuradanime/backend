package controllers

import (
	"encoding/json"
	"net/http"

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
func getInitiatorAndReceiver(r *http.Request) (string, string, error) {
	initiator := chi.URLParam(r, "initiator")
	receiver := chi.URLParam(r, "receiver")

	if initiator == "" || receiver == "" {
		return "", "", http.ErrMissingFile
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

	targetUser := chi.URLParam(r, "userID")
	if targetUser == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
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
