package app

import (
	"log"

	"github.com/afuradanime/backend/internal/adapters/controllers"
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-chi/chi/v5"
)

func (a *Application) InitRoutes(r *chi.Mux) {

	log.Printf("[Routing] Initializing public routes...")
	r.Group(func(r chi.Router) {
		r.Mount("/auth", a.BootstrapAuthModule())
		log.Println("[Routing] Auth routes initialized")
		r.Mount("/anime", a.BootstrapAnimeModule())
		log.Println("[Routing] Anime routes initialized")

	})

	log.Printf("[Routing] Initializing protected routes...")
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))
		r.Mount("/users", a.BootstrapUserModule())
		log.Println("[Routing] User routes initialized")
		r.Mount("/friends", a.BootstrapFriendsModule())
		log.Println("[Routing] Friendship routes initialized")
	})

	log.Println("[Routing] All routes initialized successfully!")
}

func (a *Application) BootstrapUserModule() chi.Router {
	userRepo := repositories.NewUserRepository(a.Mongo)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := chi.NewRouter()
	r.Get("/", userController.GetUsers)
	r.Get("/{id}", userController.GetUserByID)
	r.Put("/{id}", userController.UpdateUserInfo)

	return r
}

func (a *Application) BootstrapAnimeModule() chi.Router {
	animeRepo := repositories.NewAnimeRepository()
	animeService := services.NewAnimeService(animeRepo)
	animeController := controllers.NewAnimeController(animeService)

	r := chi.NewRouter()
	r.Get("/{id}", animeController.GetAnimeByID)
	r.Get("/search", animeController.SearchAnime)
	r.Get("/seasonal", animeController.GetAnimeThisSeason)
	r.Get("/studio/{id}", animeController.GetAnimeByStudioID)
	r.Get("/producer/{id}", animeController.GetAnimeByProducerID)
	r.Get("/licensor/{id}", animeController.GetAnimeByLicensorID)

	return r
}

func (a *Application) BootstrapFriendsModule() chi.Router {
	friendshipRepo := repositories.NewFriendshipRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipService := services.NewFriendshipService(userRepo, friendshipRepo)
	friendshipController := controllers.NewFriendshipController(friendshipService)

	r := chi.NewRouter()
	r.Put("/send/{initiator}/{receiver}", friendshipController.SendFriendRequest)
	r.Put("/accept/{initiator}/{receiver}", friendshipController.AcceptFriendRequest)
	r.Put("/decline/{initiator}/{receiver}", friendshipController.DeclineFriendRequest)
	r.Put("/block/{initiator}/{receiver}", friendshipController.BlockUser)
	r.Get("/{userID}", friendshipController.ListFriends)
	r.Get("/pending/{userID}", friendshipController.ListPendingFriendRequests)
	r.Get("/check/{initiator}/{receiver}", friendshipController.AreFriends)

	return r
}

func (a *Application) BootstrapAuthModule() chi.Router {

	jwtService := services.NewJWTService(a.JWTConfig)
	userService := services.NewUserService(repositories.NewUserRepository(a.Mongo))
	googleAuthController := controllers.NewGoogleAuthController(a.OAuth2Config, jwtService, userService)

	r := chi.NewRouter()
	r.Get("/google/login", googleAuthController.Login)
	r.Get("/google/callback", googleAuthController.Callback)
	r.Get("/me", googleAuthController.WhoAmI)
	r.Get("/logout", googleAuthController.Logout)

	return r
}
