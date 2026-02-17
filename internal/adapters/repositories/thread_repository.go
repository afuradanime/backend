package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type ThreadRepository struct {
	threadCollection  *mongo.Collection
	contextCollection *mongo.Collection
}

func NewThreadRepository(db *mongo.Database) *ThreadRepository {
	return &ThreadRepository{
		threadCollection:  db.Collection("threads"),
		contextCollection: db.Collection("thread_contexts"),
	}
}

func (t *ThreadRepository) CreateThreadContext(ctx context.Context, thcontext *domain.ThreadContext) (*domain.ThreadContext, error) {
	_, err := t.contextCollection.InsertOne(ctx, thcontext)
	if err != nil {
		return nil, err
	}
	return thcontext, nil
}

func (t *ThreadRepository) CreateThreadPost(ctx context.Context, post *domain.ThreadPost) (*domain.ThreadPost, error) {
	_, err := t.threadCollection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}
