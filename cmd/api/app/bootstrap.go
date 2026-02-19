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

	// Bootstrap a thread conversation on Kray's profile
	BootstrapThreadPosts(context.Background(), threadRepo, krayID, taikoID, testID)
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
	userTaiko.RewardBadge(value.UserBadgeSuperMegaIllyaFan)

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userTaiko)
	if err != nil {
		panic(err)
	}
	taikoID = userTaiko.ID

	userTest, err := domain.NewUser("Afuradanime", "teste@mail.teste")
	if err != nil {
		panic(err)
	}

	userTest.UpdateLocation("Afurada")
	userTest.UpdateAllowsFriendRequests(false)
	userTest.UpdateAllowsRecommendations(false)

	userTest.RewardBadge(value.UserBadgeBrand)

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

func BootstrapThreadPosts(ctx context.Context, threadRepo *repositories.ThreadRepository, krayID, taikoID, testID int) {

	post1 := domain.NewThreadPost(krayID, krayID, "Bem vindos ao meu perfil! 🎉")
	post1.CreatedAt = time.Now().Add(-48 * time.Hour).Unix()
	_, err := threadRepo.CreateThreadPost(ctx, post1)
	if err != nil {
		panic(err)
	}

	post2 := domain.NewThreadPost(krayID, taikoID, "Grande perfil! Parabéns pela criação do site 👏")
	post2.CreatedAt = time.Now().Add(-36 * time.Hour).Unix()
	post2.ReplyToPost(post1.ID)
	_, err = threadRepo.CreateThreadPost(ctx, post2)
	if err != nil {
		panic(err)
	}

	post5 := domain.NewThreadPost(krayID, testID, "Que perfil fixe ⭐")
	post5.CreatedAt = time.Now().Add(-36 * time.Hour).Unix()
	post5.ReplyToPost(post1.ID)
	_, err = threadRepo.CreateThreadPost(ctx, post5)
	if err != nil {
		panic(err)
	}

	post3 := domain.NewThreadPost(krayID, krayID, "Obrigado! Ainda há muito para fazer 😄")
	post3.CreatedAt = time.Now().Add(-24 * time.Hour).Unix()
	post3.ReplyToPost(post2.ID)
	_, err = threadRepo.CreateThreadPost(ctx, post3)
	if err != nil {
		panic(err)
	}

	post4 := domain.NewThreadPost(krayID, testID, "Olá, acabei de me juntar! 🙋")
	post4.CreatedAt = time.Now().Add(-12 * time.Hour).Unix()
	_, err = threadRepo.CreateThreadPost(ctx, post4)
	if err != nil {
		panic(err)
	}

	post6 := domain.NewThreadPost(krayID, taikoID, "Bem vindo ao site! 🎊")
	post6.CreatedAt = time.Now().Add(-6 * time.Hour).Unix()
	post6.ReplyToPost(post4.ID)
	_, err = threadRepo.CreateThreadPost(ctx, post6)
	if err != nil {
		panic(err)
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
