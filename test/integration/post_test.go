package integration

import (
	"context"
	"strconv"
	"testing"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/services"
	tests "github.com/afuradanime/backend/test"
	"github.com/stretchr/testify/require"
)

func TestSendPost(t *testing.T) {

	USER1 := 1

	app, cleanup := tests.SetupTestApp(t)
	defer cleanup()

	ctx := context.Background()

	postRepo := repositories.NewPostRepository(app.Mongo)
	userRepo := repositories.NewUserRepository(app.Mongo)

	animeSrv := services.NewAnimeService(repositories.NewAnimeRepository())

	groupRepo := repositories.NewGroupRepository(app.Mongo)
	groupServ := services.NewGroupService(groupRepo, userRepo)

	friendRepo := repositories.NewFriendshipRepository(app.Mongo)
	friendServ := services.NewFriendshipService(userRepo, friendRepo)

	service := services.NewPostService(postRepo, userRepo, friendServ, animeSrv, groupServ)

	p, err := service.CreatePost(ctx, strconv.Itoa(USER1), value.ParentTypeUser, "Test post", USER1)
	require.NoError(t, err)
	require.NotNil(t, p)
}
