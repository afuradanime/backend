package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
)

type ThreadsRepository interface {
	CreateThreadContext(ctx context.Context, context *domain.ThreadContext) (*domain.ThreadContext, error)
	CreateThreadPost(ctx context.Context, post *domain.ThreadPost) (*domain.ThreadPost, error)
	GetThreadContextByID(ctx context.Context, id int) (*domain.ThreadContext, error)
	GetThreadPostsByContext(ctx context.Context, contextId int) ([]*domain.ThreadPost, error)
}

type ThreadsService interface {
	CreateThreadPost(ctx context.Context, contextId int, contextType string, posterId int, content string) (*domain.ThreadPost, error)
	GetThreadPostsByContext(ctx context.Context, contextId int) ([]*domain.ThreadPost, error)
}
