package repositories

import (
	"context"
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnimeListRepository struct {
	collection *mongo.Collection
}

func NewAnimeListRepository(db *mongo.Database) *AnimeListRepository {
	return &AnimeListRepository{
		collection: db.Collection("anime_lists"),
	}
}

func (r *AnimeListRepository) FetchUserList(ctx context.Context, userID int) (*domain.UserAnimeList, error) {
	filter := bson.M{"user_id": userID}

	var list domain.UserAnimeList
	err := r.collection.FindOne(ctx, filter).Decode(&list)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &list, nil
}

func (r *AnimeListRepository) SaveUserList(ctx context.Context, list *domain.UserAnimeList) error {
	filter := bson.M{"user_id": list.UserID}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(ctx, filter, list, opts)
	return err
}
