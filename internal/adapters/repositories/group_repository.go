package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupRepository struct {
	collection *mongo.Collection
}

func NewGroupRepository(db *mongo.Database) *GroupRepository {
	return &GroupRepository{
		collection: db.Collection("groups"),
	}
}

func (r *GroupRepository) CreateGroup(ctx context.Context, group *domain.Group) error {
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

func (r *GroupRepository) GetGroup(ctx context.Context, groupId string) (*domain.Group, error) {
	var group domain.Group
	err := r.collection.FindOne(ctx, bson.M{
		"_id": groupId,
	}).Decode(&group)

	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) GetGroups(ctx context.Context, pageNumber, pageSize int) ([]*domain.Group, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)),
	)
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	defer cursor.Close(ctx)

	var groups []*domain.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	return groups, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, group *domain.Group) error {
	filter := bson.M{"_id": group.ID}
	update := bson.M{"$set": group}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
