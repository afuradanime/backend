package repositories

import (
	"context"
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingCacheRepository struct {
	collection *mongo.Collection
}

func NewRatingCacheRepository(db *mongo.Database) *RatingCacheRepository {
	return &RatingCacheRepository{
		collection: db.Collection("rating_cache"),
	}
}

func (r *RatingCacheRepository) CreateRatingCache(ctx context.Context, cache *domain.RatingCache) error {
	_, err := r.collection.InsertOne(ctx, cache)
	return err
}

func (r *RatingCacheRepository) UpdateRatingCache(ctx context.Context, cache *domain.RatingCache) error {
	filter := bson.M{"anime_id": cache.AnimeID}
	update := bson.M{"$set": cache}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *RatingCacheRepository) RemoveRating(ctx context.Context, userID int, animeID int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"anime_id": animeID})
	return err
}

func (r *RatingCacheRepository) GetRatingCache(ctx context.Context, animeID int) (*domain.RatingCache, error) {
	filter := bson.M{"anime_id": animeID}

	var cache domain.RatingCache
	err := r.collection.FindOne(ctx, filter).Decode(&cache)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &cache, nil
}
