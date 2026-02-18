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

type FriendshipRepository struct {
	collection *mongo.Collection
}

func NewFriendshipRepository(db *mongo.Database) *FriendshipRepository {
	return &FriendshipRepository{
		collection: db.Collection("friendships"),
	}
}

func (r *FriendshipRepository) CreateFriendship(ctx context.Context, friendship *domain.Friendship) error {
	_, err := r.collection.InsertOne(ctx, friendship)
	return err
}

func (r *FriendshipRepository) GetFriendship(ctx context.Context, initiator int, receiver int) (*domain.Friendship, error) {
	var friendship domain.Friendship
	err := r.collection.FindOne(ctx, bson.M{
		"initiator": initiator,
		"receiver":  receiver,
	}).Decode(&friendship)

	if err != nil {
		return nil, err
	}

	return &friendship, nil
}

func (r *FriendshipRepository) UpdateFriendshipStatus(ctx context.Context, initiator int, receiver int, status value.FriendshipStatus) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"initiator": initiator,
		"receiver":  receiver,
	}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})

	return err
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, userId int, pageNumber, pageSize int) ([]int, utils.Pagination, error) {

	skip := (pageNumber - 1) * pageSize

	filter := bson.M{
		"$or": []bson.M{
			{"initiator": userId},
			{"receiver": userId},
		},
		"status": value.FriendshipStatusAccepted,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	// Get friendships where user is initiator or receiver and status is accepted
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var friendships []domain.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, utils.Pagination{}, err
	}

	// Extract friend IDs
	friendIds := make([]int, len(friendships))
	for i, f := range friendships {
		if f.Initiator == userId {
			friendIds[i] = f.Receiver
		} else {
			friendIds[i] = f.Initiator
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return friendIds, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *FriendshipRepository) GetPendingFriendRequests(
	ctx context.Context,
	userId int,
	pageNumber, pageSize int,
) ([]int, utils.Pagination, error) {

	skip := (pageNumber - 1) * pageSize

	filter := bson.M{
		"receiver": userId,
		"status":   value.FriendshipStatusPending,
	}

	// Count total pending requests
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"_id": -1}) // newest requests first

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var friendships []domain.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, utils.Pagination{}, err
	}

	// Extract initiator IDs
	requestIds := make([]int, len(friendships))
	for i, f := range friendships {
		requestIds[i] = f.Initiator
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return requestIds, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
