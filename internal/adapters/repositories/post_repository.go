package repositories

import (
	"context"
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *PostRepository {
	return &PostRepository{
		collection: db.Collection("posts"),
	}
}

func (r *PostRepository) GetPostById(ctx context.Context, postID string) (*domain.Post, error) {
	var post domain.Post
	err := r.collection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetPostReplies(ctx context.Context, parentID string) ([]*domain.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"parent_id": parentID})
	if err != nil {
		return nil, errors.New("failed to fetch post replies: " + err.Error())
	}

	var replies []*domain.Post
	for cursor.Next(ctx) {
		var reply domain.Post
		if err := cursor.Decode(&reply); err != nil {
			return nil, err
		}
		replies = append(replies, &reply)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return replies, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	_, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	// Used for updates and soft deletes (deletes that keep the post but mark it as deleted by nullifying the text and createdBy fields)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": post.ID}, bson.M{"$set": post})
	return err
}

func (r *PostRepository) AddReplyToPost(ctx context.Context, parentPostID string, replyID string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": parentPostID}, bson.M{"$push": bson.M{"posts": replyID}})
	return err
}
