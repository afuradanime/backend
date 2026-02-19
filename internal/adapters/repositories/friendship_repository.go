package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *FriendshipRepository) UpdateFriendship(ctx context.Context, f *domain.Friendship) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"initiator": f.Initiator,
		"receiver":  f.Receiver,
	}, bson.M{
		"$set": bson.M{
			"status": f.Status,
		},
	})
	return err
}

func (r *FriendshipRepository) DeleteFriendship(ctx context.Context, initiator int, receiver int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"$or": []bson.M{
			{
				"initiator": initiator,
				"receiver":  receiver,
			},
			{
				"initiator": receiver,
				"receiver":  initiator,
			},
		},
	})
	return err
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, userId int, pageNumber, pageSize int) ([]domain.User, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize

	// Get the friendship list for this user
	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"$or": []bson.M{
			{"initiator": userId},
			{"receiver": userId},
		},
		"status": value.FriendshipStatusAccepted,
	}}}

	// Chain a lookup to not query twice
	lookupInitiator := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "initiator",
		"foreignField": "_id",
		"as":           "initiator_user",
	}}}

	lookupReceiver := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "receiver",
		"foreignField": "_id",
		"as":           "receiver_user",
	}}}

	// Pick the friend (the one that isn't the requesting user)
	// Get the friend field based on who initiated,
	// If initiator = userId, project receiver_user, else project initiator_user
	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.M{
		"friend": bson.M{
			"$cond": bson.M{
				"if":   bson.M{"$eq": []interface{}{"$initiator", userId}},
				"then": bson.M{"$arrayElemAt": []interface{}{"$receiver_user", 0}},
				"else": bson.M{"$arrayElemAt": []interface{}{"$initiator_user", 0}},
			},
		},
	}}}

	// Replace pipeline result with a user document, we don't need the rest
	replaceRootStage := bson.D{{Key: "$replaceRoot", Value: bson.M{
		"newRoot": "$friend",
	}}}

	// Count pipeline
	countPipeline := mongo.Pipeline{matchStage}
	countStage := bson.D{{Key: "$count", Value: "total"}}
	countCursor, err := r.collection.Aggregate(ctx, append(countPipeline, countStage))

	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var countResult []struct {
		Total int `bson:"total"`
	}

	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, utils.Pagination{}, err
	}

	total := 0
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Main pipeline
	skipStage := bson.D{{Key: "$skip", Value: int64(skip)}}
	limitStage := bson.D{{Key: "$limit", Value: int64(pageSize)}}

	pipeline := mongo.Pipeline{
		matchStage,
		lookupInitiator,
		lookupReceiver,
		addFieldsStage,
		replaceRootStage,
		skipStage,
		limitStage,
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return users, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *FriendshipRepository) GetPendingFriendRequests(ctx context.Context, userId int, pageNumber, pageSize int) ([]domain.User, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize

	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"receiver": userId,
		"status":   value.FriendshipStatusPending,
	}}}

	lookupInitiator := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "initiator",
		"foreignField": "_id",
		"as":           "initiator_user",
	}}}

	replaceRootStage := bson.D{{Key: "$replaceRoot", Value: bson.M{
		"newRoot": bson.M{"$arrayElemAt": []interface{}{"$initiator_user", 0}},
	}}}

	// Count
	countCursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		bson.D{{Key: "$count", Value: "total"}},
	})

	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var countResult []struct {
		Total int `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, utils.Pagination{}, err
	}

	total := 0
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Main pipeline
	pipeline := mongo.Pipeline{
		matchStage,
		lookupInitiator,
		replaceRootStage,
		bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}},
		bson.D{{Key: "$skip", Value: int64(skip)}},
		bson.D{{Key: "$limit", Value: int64(pageSize)}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return users, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
