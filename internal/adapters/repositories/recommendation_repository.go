package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecommendationRepository struct {
	collection *mongo.Collection
}

func NewRecommendationRepository(db *mongo.Database) *RecommendationRepository {
	return &RecommendationRepository{
		collection: db.Collection("recommendations"),
	}
}

func (r *RecommendationRepository) Create(ctx context.Context, rec *domain.Recommendation) error {

	_, err := r.collection.InsertOne(ctx, rec)
	return err
}

func (r *RecommendationRepository) HasBeenRecommended(ctx context.Context, receiverID, animeID int) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"receiver": receiverID,
		"anime":    animeID,
	})
	return count > 0, err
}

func (r *RecommendationRepository) RecommendationStackCount(ctx context.Context, receiverID int) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"receiver": receiverID,
	})
	return count, err
}

func (r *RecommendationRepository) GetForUser(ctx context.Context, receiverID, pageNumber, pageSize int) ([]*domain.Recommendation, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize

	filter := bson.M{"receiver": receiverID}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	defer cursor.Close(ctx)

	var recs []*domain.Recommendation
	if err := cursor.All(ctx, &recs); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	return recs, utils.Pagination{PageNumber: pageNumber, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (r *RecommendationRepository) DismissRecommendation(ctx context.Context, receiverID, anime int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"receiver": receiverID,
		"anime":    anime,
	})
	return err
}
