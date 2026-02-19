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

func (t *ThreadRepository) GetThreadContextByID(ctx context.Context, id int) (*domain.ThreadContext, error) {
	var thcontext domain.ThreadContext
	err := t.contextCollection.FindOne(
		ctx,
		map[string]interface{}{"contextId": id},
	).Decode(&thcontext)
	if err != nil {
		return nil, err
	}
	return &thcontext, nil
}

func (t *ThreadRepository) GetThreadPostsByContext(ctx context.Context, contextId int) ([]*domain.ThreadPost, error) {
	cursor, err := t.threadCollection.Find(
		ctx,
		map[string]interface{}{"contextId": contextId},
	)
	if err != nil {
		return nil, err
	}

	var posts []*domain.ThreadPost
	for cursor.Next(ctx) {
		var post domain.ThreadPost
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := cursor.Err(); err != nil {
		cursor.Close(ctx)
		return nil, err
	}

	cursor.Close(ctx)
	return posts, nil
}
