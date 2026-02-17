package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
)

type ThreadsRepository interface {
	CreateThreadContext(ctx context.Context, context *domain.ThreadContext) (*domain.ThreadContext, error)
	CreateThreadPost(ctx context.Context, post *domain.ThreadPost) (*domain.ThreadPost, error)
}

type ThreadsService interface {
	CreateThreadPost(ctx context.Context, context int, userId int, content string) (*domain.ThreadPost, error)
}
