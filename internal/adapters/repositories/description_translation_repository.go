package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DescriptionTranslationRepository struct {
	collection        *mongo.Collection
	counterCollection *mongo.Collection
}

func NewDescriptionTranslationRepository(db *mongo.Database) *DescriptionTranslationRepository {
	return &DescriptionTranslationRepository{
		collection:        db.Collection("description_translations"),
		counterCollection: db.Collection("counters"),
	}
}

func (r *DescriptionTranslationRepository) getNextSequence(ctx context.Context, name string) (int, error) {
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var result Counter
	err := r.counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Seq, nil
}

func (r *DescriptionTranslationRepository) CreateTranslation(ctx context.Context, translation *domain.DescriptionTranslation) error {
	nextID, err := r.getNextSequence(ctx, "description_translation_id")
	if err != nil {
		return err
	}
	translation.ID = nextID
	_, err = r.collection.InsertOne(ctx, translation)
	return err
}

func (r *DescriptionTranslationRepository) GetTranslationByID(ctx context.Context, id int) (*domain.DescriptionTranslation, error) {
	var translation domain.DescriptionTranslation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&translation)
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *DescriptionTranslationRepository) GetTranslationByAnime(ctx context.Context, anime int) (*domain.DescriptionTranslation, error) {
	var translation domain.DescriptionTranslation
	err := r.collection.FindOne(ctx, bson.M{
		"anime":  anime,
		"status": value.DescriptionTranslationApproved,
	}).Decode(&translation)
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *DescriptionTranslationRepository) GetTranslationByAnimeFromUser(ctx context.Context, anime int, id int) (*domain.DescriptionTranslation, error) {
	var translation domain.DescriptionTranslation
	err := r.collection.FindOne(ctx, bson.M{
		"anime":      anime,
		"created_by": id,
	}).Decode(&translation)
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *DescriptionTranslationRepository) GetTranslationsByUser(
	ctx context.Context,
	userID int,
	pageNumber, pageSize int,
) ([]domain.DescriptionTranslation, utils.Pagination, error) {

	skip := (pageNumber - 1) * pageSize

	filter := bson.M{"created_by": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"_id": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var translations []domain.DescriptionTranslation
	if err := cursor.All(ctx, &translations); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return translations, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *DescriptionTranslationRepository) GetPendingTranslations(
	ctx context.Context,
	pageNumber, pageSize int,
) ([]domain.DescriptionTranslation, utils.Pagination, error) {

	skip := (pageNumber - 1) * pageSize

	filter := bson.M{"status": value.DescriptionTranslationPending}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"_id": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var translations []domain.DescriptionTranslation
	if err := cursor.All(ctx, &translations); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return translations, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *DescriptionTranslationRepository) UpdateTranslation(ctx context.Context, t *domain.DescriptionTranslation) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{
		"$set": bson.M{
			"status":      t.TranslationStatus,
			"accepted_by": t.AcceptedBy,
			"accepted_at": t.AcceptedAt,
		},
	})
	return err
}

func (r *DescriptionTranslationRepository) DeleteTranslation(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
