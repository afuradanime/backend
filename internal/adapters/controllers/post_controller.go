package controllers

import (
	"errors"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/go-fuego/fuego"
)

type PostController struct {
	postService interfaces.PostService
}

func NewPostController(postService interfaces.PostService) *PostController {
	return &PostController{postService: postService}
}

func (c *PostController) GetPostById(ctx fuego.ContextNoBody) (*domain.Post, error) {
	postId := ctx.PathParam("post_id")

	post, err := c.postService.GetPostById(ctx.Context(), postId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		return nil, fuego.NotFoundError{Detail: err.Error()}
	} else if err != nil {
		return nil, fuego.InternalServerError{Detail: "Internal error when fetching post: " + err.Error()}
	}

	return post, nil
}

type CreatePostBody struct {
	Text       string               `json:"text"`
	ParentID   string               `json:"parent_id"`
	ParentType value.PostParentType `json:"parent_type"`
}

func (c *PostController) CreatePost(ctx fuego.ContextWithBody[CreatePostBody]) (*domain.Post, error) {
	posterId, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.ForbiddenError{Detail: "Not logged in, cannot create post"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Failed to decode request body: " + err.Error()}
	}

	post, err := c.postService.CreatePost(ctx.Context(), body.ParentID, body.ParentType, body.Text, posterId)
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Failed to create post: " + err.Error()}
	}

	return post, nil
}

func (c *PostController) DeletePost(ctx fuego.ContextNoBody) (any, error) {
	postId := ctx.PathParam("post_id")

	deleterId, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.ForbiddenError{Detail: "Not logged in, cannot delete post"}
	}

	err := c.postService.DeletePost(ctx.Context(), postId, deleterId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		return nil, fuego.NotFoundError{Detail: err.Error()}
	} else if errors.Is(err, domain_errors.PostDeletedError{}) {
		return nil, fuego.HTTPError{Status: 410, Detail: "Trying to delete a deleted post: " + err.Error()}
	} else if errors.Is(err, domain_errors.NotPostOwnerError{UserID: strconv.Itoa(deleterId), PostID: postId}) {
		return nil, fuego.UnauthorizedError{Detail: err.Error()}
	} else if err != nil {
		return nil, fuego.InternalServerError{Detail: "Internal error when deleting post: " + err.Error()}
	}

	return nil, nil
}

func (c *PostController) GetPostReplies(ctx fuego.ContextNoBody) ([]*domain.Post, error) {
	parentId := ctx.PathParam("parent_id")

	replies, err := c.postService.GetPostReplies(ctx.Context(), parentId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		return nil, fuego.NotFoundError{Detail: err.Error()}
	} else if err != nil {
		return nil, fuego.InternalServerError{Detail: "Internal error when fetching post replies: " + err.Error()}
	}

	return replies, nil
}

type CreateReplyBody struct {
	Text string `json:"text"`
}

func (c *PostController) CreateReply(ctx fuego.ContextWithBody[CreateReplyBody]) (*domain.Post, error) {
	posterId, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.ForbiddenError{Detail: "Not logged in, cannot create post"}
	}

	replyToPostID := ctx.PathParam("post_id")
	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Failed to decode request body: " + err.Error()}
	}

	post, err := c.postService.CreateReply(ctx.Context(), replyToPostID, body.Text, posterId)
	if errors.Is(err, domain_errors.PostNotFoundError{}) {
		return nil, fuego.NotFoundError{Detail: err.Error()}
	} else if err != nil {
		return nil, fuego.BadRequestError{Detail: "Failed to create reply: " + err.Error()}
	}

	return post, nil
}
