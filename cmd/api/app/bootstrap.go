package app

import (
	"context"

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

	userKray, err := domain.NewUser("1", "kray", "kray@afurada.anime")
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

	userTaiko, err := domain.NewUser("2", "taiko", "taiko@afurada.anime")
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

	_, err = userCollection.InsertMany(
		ctx,
		[]interface{}{
			userKray,
			userTaiko,
		},
	)

	if err != nil {
		panic(err)
	}
}

func BootstrapFriendships(ctx context.Context, friendshipCollection *mongo.Collection) {

	friendship := domain.NewFriendRequest("1", "2")
	friendship.Accept()

	_, err := friendshipCollection.InsertOne(ctx, friendship)
	if err != nil {
		panic(err)
	}
}
