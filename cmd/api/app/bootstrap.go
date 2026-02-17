package app

import (
	"context"
	"time"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
)

func (a *Application) Bootstrap() {

	// Nuke EVERYTHING
	a.Mongo.Drop(context.Background())

	// Create repositories
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipRepo := repositories.NewFriendshipRepository(a.Mongo)

	// Bootstrap users and get their auto-generated IDs
	krayID, taikoID, testID := BootstrapUsers(context.Background(), userRepo)

	// Bootstrap thread contexts for user profiles
	threadRepo := repositories.NewThreadRepository(a.Mongo)
	BootstrapThreadContexts(context.Background(), threadRepo, krayID, taikoID, testID)

	// Bootstrap friendships using the actual user IDs
	BootstrapFriendships(context.Background(), friendshipRepo, krayID, taikoID, testID)
}

func BootstrapUsers(ctx context.Context, userRepo *repositories.UserRepository) (krayID, taikoID, testID int) {

	userKray, err := domain.NewUser("KrayRui", "kray@afurada.anime")
	if err != nil {
		panic(err)
	}

	userKray.UpdateLocation("Porto")
	userKray.UpdateSocials([]string{
		"https://x.com/RuiIshigami",
		"https://github.com/Rui-San",
	})

	userKray.AddRole(value.UserRoleModerator)
	userKray.AddRole(value.UserRoleAdmin)

	userKray.UpdateAvatarURL("/pfps/d7dea5d3e09941f563dabf364b4db31cac63a5f1.png")

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userKray)
	if err != nil {
		panic(err)
	}
	krayID = userKray.ID

	userTaiko, err := domain.NewUser("Sagiri719", "taiko@afurada.anime")
	if err != nil {
		panic(err)
	}

	userTaiko.UpdateLocation("Porto")
	userTaiko.UpdateSocials([]string{
		"https://x.com/Sagiri719",
		"https://github.com/Sagiri721",
	})

	userTaiko.AddRole(value.UserRoleModerator)
	userTaiko.AddRole(value.UserRoleAdmin)

	userTaiko.UpdateAvatarURL("/pfps/e59084c01caf44df3c240a3c78009d080ea02556.png")

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userTaiko)
	if err != nil {
		panic(err)
	}
	taikoID = userTaiko.ID

	userTest, err := domain.NewUser("Teste", "teste@mail.teste")
	if err != nil {
		panic(err)
	}

	userTest.UpdateLocation("Porto")
	userTest.UpdateSocials([]string{
		"https://x.com/Teste",
	})
	userTest.UpdateBirthday(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	userTest.UpdatePronouns("user/teste")

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userTest)
	if err != nil {
		panic(err)
	}
	testID = userTest.ID

	return krayID, taikoID, testID
}

func BootstrapThreadContexts(ctx context.Context, threadRepo *repositories.ThreadRepository, krayID, taikoID, testID int) {
	for _, userID := range []int{krayID, taikoID, testID} {
		tc := domain.NewContext(userID, "Profile")
		_, err := threadRepo.CreateThreadContext(ctx, tc)
		if err != nil {
			panic(err)
		}
	}
}

func BootstrapFriendships(ctx context.Context, friendshipRepo *repositories.FriendshipRepository, krayID, taikoID, testID int) {

	friendship := domain.NewFriendRequest(krayID, taikoID)
	friendship.Accept()

	err := friendshipRepo.CreateFriendship(ctx, friendship)
	if err != nil {
		panic(err)
	}

	friendship2 := domain.NewFriendRequest(taikoID, testID)
	friendship2.Accept()

	err = friendshipRepo.CreateFriendship(ctx, friendship2)
	if err != nil {
		panic(err)
	}
}
