package repositories

import (
	"context"
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) UpdatePersonalInfo(
	ctx context.Context,
	id string,
	email *string,
	username *string,
	location *string,
	pronouns *string,
	socials *[]string,
) error {

	setFields := bson.M{}

	if email != nil {
		setFields["email"] = *email
	}

	if username != nil {
		setFields["username"] = *username
	}

	if location != nil {
		setFields["location"] = *location
	}

	if pronouns != nil {
		setFields["pronouns"] = *pronouns
	}

	if socials != nil {
		setFields["socials"] = *socials
	}

	if len(setFields) == 0 {
		return nil // nothing to update
	}

	update := bson.M{"$set": setFields}

	result, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
