package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/afuradanime/backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Counter struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

type UserRepository struct {
	collection        *mongo.Collection
	counterCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection:        db.Collection("users"),
		counterCollection: db.Collection("counters"),
	}
}

func (r *UserRepository) getNextSequence(ctx context.Context, name string) (int, error) {
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

func (r *UserRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}

	return &user, nil
}

func (r *UserRepository) GetUserByProvider(ctx context.Context, provider string, providerID string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(
		ctx, bson.M{"provider": provider, "provider_id": providerID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Get next auto-incrementing ID
	nextID, err := r.getNextSequence(ctx, "user_id")
	if err != nil {
		return err
	}

	user.ID = nextID
	user.CreatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) UpdatePersonalInfo(
	ctx context.Context,
	id int,
	user *domain.User,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"email":    user.Email,
			"username": user.Username,
			"location": user.Location,
			"pronouns": user.Pronouns,
			"socials":  user.Socials,
		}},
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"last_login": time.Now(),
		}},
	)
	return err
}
