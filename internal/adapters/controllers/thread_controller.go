package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/go-chi/chi/v5"
)

type ThreadController struct {
	threadService interfaces.ThreadsService
}

func NewThreadController(threadService interfaces.ThreadsService) *ThreadController {
	return &ThreadController{
		threadService: threadService,
	}
}

func (c *ThreadController) CreateThreadPost(w http.ResponseWriter, r *http.Request) {
	contextId, err := strconv.Atoi(chi.URLParam(r, "contextId"))
	if err != nil {
		http.Error(w, "Invalid context ID", http.StatusBadRequest)
		return
	}

	contextType := chi.URLParam(r, "contextType")
	if contextType == "" {
		http.Error(w, "Context type is required", http.StatusBadRequest)
		return
	}

	posterId, ok := r.Context().Value(middlewares.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Content string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if body.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	post, err := c.threadService.CreateThreadPost(r.Context(), contextId, contextType, posterId, body.Content)
	if err != nil {
		http.Error(w, "Failed to create thread post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (c *ThreadController) GetThreadPostsByContext(w http.ResponseWriter, r *http.Request) {
	contextId, err := strconv.Atoi(chi.URLParam(r, "contextId"))
	if err != nil {
		http.Error(w, "Invalid context ID", http.StatusBadRequest)
		return
	}

	posts, err := c.threadService.GetThreadPostsByContext(r.Context(), contextId)
	if err != nil {
		http.Error(w, "Failed to get thread posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
