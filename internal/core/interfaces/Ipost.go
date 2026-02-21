package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
)

type PostRepository interface {
	GetPostById(ctx context.Context, postID string) (*domain.Post, error)
	GetPostReplies(ctx context.Context, parentId string) ([]*domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) error
	AddReplyToPost(ctx context.Context, parentPostID string, replyID string) error
}

type PostService interface {
	GetPostById(ctx context.Context, postID string) (*domain.Post, error)
	GetPostReplies(ctx context.Context, parentId string) ([]*domain.Post, error)
	CreatePost(ctx context.Context, parentId string, parentType value.PostParentType, text string, posterId int) (*domain.Post, error)
	CreateReply(ctx context.Context, replyToPostID string, text string, createdBy int) (*domain.Post, error)
	DeletePost(ctx context.Context, postID string, deleterId int) error
}
