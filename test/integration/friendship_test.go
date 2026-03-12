package integration

import (
	"context"
	"testing"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/services"
	tests "github.com/afuradanime/backend/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)

	friendship, err := friendshipRepo.GetFriendship(ctx, 1, 2)
	require.NoError(t, err)
	require.NotNil(t, friendship)

	assert.Equal(t, 1, friendship.Initiator)
	assert.Equal(t, 2, friendship.Receiver)
	assert.Equal(t, value.FriendshipStatusPending, friendship.Status)
}

func TestSendAlreadyExistingFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to send the same friend request again
	err = service.SendFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestAcceptNonExistentFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.AcceptFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestRejectNonExistentFriendRequest(t *testing.T) {
	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.DeclineFriendRequest(ctx, USER2, USER1)
	require.Error(t, err)
}

func TestDoubleSendFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.SendFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestSendToBlockedUserFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	// Block USER1 by USER2
	err := service.BlockUser(ctx, USER2, USER1)
	require.NoError(t, err)

	// Try to send friend request from USER1 to USER2
	err = service.SendFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestAcceptFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.AcceptFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	friendship, err := friendshipRepo.GetFriendship(ctx, 1, 2)
	require.NoError(t, err)
	require.NotNil(t, friendship)

	assert.Equal(t, value.FriendshipStatusAccepted, friendship.Status)
}

func TestAcceptMyOwnFriendRequest(t *testing.T) {

	USER1 := 1

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.AcceptFriendRequest(ctx, USER1, USER1)
	require.Error(t, err)
}

func TestAcceptAlreadyAcceptedFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.AcceptFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to accept the same friend request again
	err = service.AcceptFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestAcceptFriendRequestByNonReceiver(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to accept the friend request by the initiator instead of the receiver
	err = service.AcceptFriendRequest(ctx, USER1, USER1)
	require.Error(t, err)
}

func TestRejectFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.DeclineFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	friendship, err := friendshipRepo.GetFriendship(ctx, 1, 2)
	assert.Nil(t, friendship)
}

func TestRejectAlreadyDeclinedFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.DeclineFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to decline the same friend request again
	err = service.DeclineFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestRejectFriendRequestByNonReceiver(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to reject the friend request by the initiator instead of the receiver
	err = service.DeclineFriendRequest(ctx, USER2, USER1)
	require.Error(t, err)
}

func TestBlockUser(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.BlockUser(ctx, USER1, USER2)
	require.NoError(t, err)

	blocked, err := friendshipRepo.GetFriendship(ctx, USER1, USER2)
	require.NoError(t, err)
	require.NotNil(t, blocked)

	assert.Equal(t, value.FriendshipStatusBlocked, blocked.Status)
}

func TestBlockAlreadyBlockedUser(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.BlockUser(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to block the same user again
	err = service.BlockUser(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestBlockMyself(t *testing.T) {

	USER1 := 1

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.BlockUser(ctx, USER1, USER1)
	require.Error(t, err)
}

func TestBlockUserThenSendFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.BlockUser(ctx, USER1, USER2)
	require.NoError(t, err)

	// Try to send friend request from USER2 to USER1 after being blocked
	err = service.SendFriendRequest(ctx, USER2, USER1)
	require.Error(t, err)
}

func TestBlockUserThenAcceptFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.BlockUser(ctx, USER2, USER1)
	require.NoError(t, err)

	// Try to accept the friend request after being blocked
	err = service.AcceptFriendRequest(ctx, USER1, USER2)
	require.Error(t, err)
}

func TestBlockUserThenRejectFriendRequest(t *testing.T) {

	USER1 := 1
	USER2 := 2

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER2)
	require.NoError(t, err)

	err = service.BlockUser(ctx, USER2, USER1)
	require.NoError(t, err)

	// Try to reject the friend request after being blocked
	err = service.DeclineFriendRequest(ctx, USER2, USER1)
	require.Error(t, err)
}

func TestFriendMyself(t *testing.T) {

	USER1 := 1

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	// Clean friendship collection before test
	app.Mongo.Collection("friendships").Drop(context.Background())

	ctx := context.Background()

	friendshipRepo := repositories.NewFriendshipRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	service := services.NewFriendshipService(userRepo, friendshipRepo)

	err := service.SendFriendRequest(ctx, USER1, USER1)
	require.Error(t, err)
}
