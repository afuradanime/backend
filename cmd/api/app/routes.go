package app

import (
	"net/http"

	"github.com/afuradanime/backend/internal/adapters/controllers"
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-chi/chi/v5"
)

// TODO: cors
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:5173"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Application) InitRoutes(r *chi.Mux) {
	r.Use(CORSMiddleware)

	r.Group(func(r chi.Router) {
		r.Mount("/auth", a.BootstrapAuthModule())
		r.Mount("/anime", a.BootstrapAnimeModule())
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))
		r.Mount("/users", a.BootstrapUserModule())
		r.Mount("/friends", a.BootstrapFriendsModule())
		r.Mount("/threads", a.BootstrapThreadModule())
		r.Mount("/translations", a.BootstrapTranslationsModule())
	})
}

func (a *Application) BootstrapTranslationsModule() chi.Router {
	translationRepo := repositories.NewDescriptionTranslationRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	animeRepo := repositories.NewAnimeRepository()
	translationService := services.NewDescriptionTranslationService(translationRepo, animeRepo, userRepo)
	translationController := controllers.NewDescriptionTranslationController(translationService)

	r := chi.NewRouter()

	// User routes
	r.Post("/anime/{animeID}", translationController.SubmitTranslation)
	r.Get("/anime/{animeID}", translationController.GetAnimeTranslation)
	r.Get("/me", translationController.GetMyTranslations)

	// Moderation routes
	r.Put("/{id}/accept", translationController.AcceptTranslation)
	r.Put("/{id}/reject", translationController.RejectTranslation)
	r.Get("/pending", translationController.GetPendingTranslations)

	return r
}

func (a *Application) BootstrapUserModule() chi.Router {
	userRepo := repositories.NewUserRepository(a.Mongo)
	thRepo := repositories.NewThreadRepository(a.Mongo)
	userService := services.NewUserService(userRepo, thRepo)
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
	r.Put("/send/{receiver}", friendshipController.SendFriendRequest)
	r.Put("/accept/{initiator}", friendshipController.AcceptFriendRequest)
	r.Put("/decline/{initiator}", friendshipController.DeclineFriendRequest)
	r.Put("/block/{receiver}", friendshipController.BlockUser)
	r.Get("/{userID}", friendshipController.ListFriends)
	r.Get("/pending", friendshipController.ListPendingFriendRequests)
	r.Get("/check/{receiver}", friendshipController.FetchFriendshipStatus)

	return r
}

func (a *Application) BootstrapThreadModule() chi.Router {
	threadRepo := repositories.NewThreadRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipRepo := repositories.NewFriendshipRepository(a.Mongo)
	threadService := services.NewThreadService(threadRepo, userRepo, friendshipRepo)
	threadController := controllers.NewThreadController(threadService)

	r := chi.NewRouter()
	r.Get("/{contextId}", threadController.GetThreadPostsByContext)
	r.Post("/{contextId}/{contextType}", threadController.CreateThreadPost)

	return r
}

func (a *Application) BootstrapAuthModule() chi.Router {

	jwtService := services.NewJWTService(a.JWTConfig)
	thRepo := repositories.NewThreadRepository(a.Mongo)
	userService := services.NewUserService(repositories.NewUserRepository(a.Mongo), thRepo)
	googleAuthController := controllers.NewGoogleAuthController(a.OAuth2Config, jwtService, userService)

	r := chi.NewRouter()
	r.Get("/google/login", googleAuthController.Login)
	r.Get("/google/callback", googleAuthController.Callback)
	r.Get("/me", googleAuthController.WhoAmI)
	r.Get("/logout", googleAuthController.Logout)

	return r
}
