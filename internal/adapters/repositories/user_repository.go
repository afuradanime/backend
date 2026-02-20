package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
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

func (r *UserRepository) GetUsers(ctx context.Context, pageNumber, pageSize int) ([]*domain.User, utils.Pagination, error) {
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

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	return users, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *UserRepository) SearchByUsername(ctx context.Context, username string, pageNumber, pageSize int) ([]*domain.User, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize

	// Case-insensitive partial match
	filter := bson.M{
		"username": bson.M{
			"$regex":   username,
			"$options": "i",
		},
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"_id": 1}),
	)
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	return users, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
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

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Get next auto-incrementing ID
	nextID, err := r.getNextSequence(ctx, "user_id")
	if err != nil {
		return nil, err
	}

	user.ID = nextID
	user.CreatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, user)
	return user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"email":                  user.Email,
			"username":               user.Username,
			"avatar_url":             user.AvatarURL,
			"location":               user.Location,
			"birthday":               user.Birthday,
			"pronouns":               user.Pronouns,
			"socials":                user.Socials,
			"allows_friend_requests": user.AllowsFriendRequests,
			"allows_recommendations": user.AllowsRecommendations,
			"can_post":               user.CanPost,
			"can_translate":          user.CanTranslate,
			"roles":                  user.Roles,
			"badges":                 user.Badges,
			"last_login":             user.LastLogin,
		}},
	)
	return err
}
