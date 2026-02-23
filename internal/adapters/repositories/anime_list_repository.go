package repositories

import (
	"context"
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnimeListRepository struct {
	collection *mongo.Collection
}

func NewAnimeListRepository(db *mongo.Database) *AnimeListRepository {
	return &AnimeListRepository{
		collection: db.Collection("anime_list_entries"),
	}
}

func (r *AnimeListRepository) AddListItem(ctx context.Context, item *domain.AnimeListItem) error {
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *AnimeListRepository) UpdateListItem(ctx context.Context, item *domain.AnimeListItem) error {
	filter := bson.M{
		"user_id":  item.UserID,
		"anime_id": item.AnimeID,
	}
	_, err := r.collection.ReplaceOne(ctx, filter, item)
	return err
}

func (r *AnimeListRepository) DeleteListItem(ctx context.Context, userID int, animeID uint32) error {
	filter := bson.M{
		"user_id":  userID,
		"anime_id": animeID,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *AnimeListRepository) FetchItemByUserAndAnime(ctx context.Context, userID int, animeID uint32) (*domain.AnimeListItem, error) {
	filter := bson.M{
		"user_id":  userID,
		"anime_id": animeID,
	}

	var item domain.AnimeListItem
	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func (r *AnimeListRepository) FetchUserList(ctx context.Context, userID int, status *value.AnimeListItemStatus) ([]*domain.AnimeListItem, error) {
	filter := bson.M{"user_id": userID}

	if status != nil {
		filter["status"] = *status
	}

	opts := options.Find().SetSort(bson.M{"edited_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err // Db error, something went wrong, prolly query syntax?
	}

	defer cursor.Close(ctx)

	var items []*domain.AnimeListItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err // Error decoding documents into the AnimeListItem struct
	}

	if items == nil {
		items = make([]*domain.AnimeListItem, 0)
	}

	return items, nil
}
