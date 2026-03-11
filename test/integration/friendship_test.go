package integration

import (
	"context"
	"testing"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/services"
	tests "github.com/afuradanime/backend/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendFriendRequest(t *testing.T) {

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	user1 := &domain.User{
		ID:                   1,
		AllowsFriendRequests: true,
	}

	user2 := &domain.User{
		ID:                   2,
		AllowsFriendRequests: true,
	}

	_, err := app.Mongo.Collection("users").InsertOne(ctx, user1)
	require.NoError(t, err)

	_, err = app.Mongo.Collection("users").InsertOne(ctx, user2)
	require.NoError(t, err)

	err = service.SendFriendRequest(ctx, 1, 2)
	require.NoError(t, err)

	friendship, err := friendshipRepo.GetFriendship(ctx, 1, 2)
	require.NoError(t, err)

	assert.Equal(t, 1, friendship.Initiator)
	assert.Equal(t, 2, friendship.Receiver)
}
