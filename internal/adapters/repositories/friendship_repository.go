package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
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

func (r *FriendshipRepository) GetFriendship(ctx context.Context, initiator string, receiver string) (*domain.Friendship, error) {
	var friendship domain.Friendship
	err := r.collection.FindOne(ctx, map[string]string{
		"initiator": initiator,
		"receiver":  receiver,
	}).Decode(&friendship)

	if err != nil {
		return nil, err
	}

	return &friendship, nil
}

func (r *FriendshipRepository) UpdateFriendshipStatus(ctx context.Context, initiator string, receiver string, status value.FriendshipStatus) error {
	_, err := r.collection.UpdateOne(ctx, map[string]string{
		"initiator": initiator,
		"receiver":  receiver,
	}, map[string]interface{}{
		"$set": map[string]interface{}{
			"status": status,
		},
	})

	return err
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, userId string) ([]string, error) {

	// Get friendships where user is initiator or receiver and status is accepted
	cursor, err := r.collection.Find(ctx, map[string]interface{}{
		"$or": []map[string]string{
			{"initiator": userId},
			{"receiver": userId},
		},
		"status": value.FriendshipStatusAccepted,
	})

	if err != nil {
		return nil, err
	}

	var friendships []domain.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	// Extract friend IDs
	friendIds := make([]string, len(friendships))
	for i, f := range friendships {
		if f.Initiator == userId {
			friendIds[i] = f.Receiver
		} else {
			friendIds[i] = f.Initiator
		}
	}

	return friendIds, nil
}

func (r *FriendshipRepository) GetPendingFriendRequests(ctx context.Context, userId string) ([]string, error) {

	// Get friendships where user is receiver and status is pending
	cursor, err := r.collection.Find(ctx, map[string]interface{}{
		"receiver": userId,
		"status":   value.FriendshipStatusPending,
	})

	if err != nil {
		return nil, err
	}

	var friendships []domain.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	// Extract initiator IDs
	requestIds := make([]string, len(friendships))
	for i, f := range friendships {
		requestIds[i] = f.Initiator
	}

	return requestIds, nil
}
