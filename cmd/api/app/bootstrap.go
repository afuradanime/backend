package app

import (
	"context"
	"time"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *Application) Bootstrap() {

	// Nuke EVERYTHING
	a.Mongo.Drop(context.Background())

	userCollection := a.Mongo.Collection("users")
	BootstrapUsers(context.Background(), userCollection)

	friendshipCollection := a.Mongo.Collection("friendships")
	BootstrapFriendships(context.Background(), friendshipCollection)
}

func BootstrapUsers(ctx context.Context, userCollection *mongo.Collection) {

	userKray, err := domain.NewUser("1", "KrayRui", "kray@afurada.anime")
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

	userTaiko, err := domain.NewUser("2", "Sagiri719", "taiko@afurada.anime")
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

	userTest, err := domain.NewUser("3", "Teste", "teste@mail.teste")
	if err != nil {
		panic(err)
	}

	userTest.UpdateLocation("Porto")
	userTest.UpdateSocials([]string{
		"https://x.com/Teste",
	})
	userTest.UpdateBirthday(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	userTest.UpdatePronouns("user/teste")

	_, err = userCollection.InsertMany(
		ctx,
		[]interface{}{
			userKray,
			userTaiko,
			userTest,
		},
	)

	if err != nil {
		panic(err)
	}
}

func BootstrapFriendships(ctx context.Context, friendshipCollection *mongo.Collection) {

	friendship := domain.NewFriendRequest("1", "2")
	friendship.Accept()

	friendship2 := domain.NewFriendRequest("2", "3")
	friendship2.Accept()

	_, err := friendshipCollection.InsertMany(
		ctx,
		[]interface{}{
			friendship,
			friendship2,
		},
	)
	if err != nil {
		panic(err)
	}
}
