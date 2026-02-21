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

type PostController struct {
	postService interfaces.PostService
}

func NewPostController(postService interfaces.PostService) *PostController {
	return &PostController{postService: postService}
}

// params: post_id
func (c *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")

	post, err := c.postService.GetPostById(r.Context(), postId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal error when fetching post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// params: text, parent_id, parent_type
func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	posterId, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Not logged in, cannot create post", http.StatusForbidden)
		return
	}

	var body struct {
		Text       string               `json:"text"`
		ParentID   string               `json:"parent_id"`
		ParentType value.PostParentType `json:"parent_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	post, err := c.postService.CreatePost(r.Context(), body.ParentID, body.ParentType, body.Text, posterId)

	if err != nil {
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// params: post_id
func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")

	deleterId, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Not logged in, cannot delete post", http.StatusForbidden)
		return
	}

	err := c.postService.DeletePost(r.Context(), postId, deleterId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, domain_errors.PostDeletedError{}) {
		http.Error(w, "Trying to delete a deleted post: "+err.Error(), http.StatusGone)
		return
	} else if errors.Is(err, domain_errors.NotPostOwnerError{UserID: strconv.Itoa(deleterId), PostID: postId}) {
		// The user that tried to do this knows about our internal architecture, they shouldn't be doing this,
		// Therefore we could eventually report them internally
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Internal error when deleting post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// params: parent_id
func (c *PostController) GetPostReplies(w http.ResponseWriter, r *http.Request) {
	parentId := chi.URLParam(r, "parent_id")

	replies, err := c.postService.GetPostReplies(r.Context(), parentId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal error when fetching post replies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}

// params: post_id (path), text (body)
func (c *PostController) CreateReply(w http.ResponseWriter, r *http.Request) {
	posterId, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Not logged in, cannot create post", http.StatusForbidden)
		return
	}

	replyToPostID := chi.URLParam(r, "post_id")
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	post, err := c.postService.CreateReply(r.Context(), replyToPostID, body.Text, posterId)

	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to create reply: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
