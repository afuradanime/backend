package app

import (
	"github.com/afuradanime/backend/internal/adapters/controllers"
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Application) InitRoutes(r *chi.Mux) {

	// Global rate limiter for all regular application endpoints
	// 10 Second limiting every 30 bursts
	globalLimiter := &middlewares.IPRateLimiter{Rps: 10, Burst: 30}

	r.Use(
		middleware.Logger,
		middleware.Recoverer, // useful middleware to recover from panics and return a 500 error
		middlewares.CORSMiddleware,
		globalLimiter.Middleware,
	)

	r.Group(func(r chi.Router) {
		r.Mount("/auth", a.BootstrapAuthModule())
		r.Mount("/anime", a.BootstrapAnimeModule())

		// Partially public
		r.Mount("/users", a.BootstrapUserModule())
		r.Mount("/friends", a.BootstrapFriendsModule())
		r.Mount("/translations", a.BootstrapTranslationsModule())
		r.Mount("/animelist", a.BootstrapAnimeListModule())
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))
		r.Mount("/reports", a.BootstrapReportsModule())
		r.Mount("/posts", a.BootstrapPostModule())
		r.Mount("/recommendations", a.BootstrapRecommendationsModule())
	})
}

func (a *Application) BootstrapTranslationsModule() chi.Router {
	translationRepo := repositories.NewDescriptionTranslationRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	animeRepo := repositories.NewAnimeRepository()
	translationService := services.NewDescriptionTranslationService(translationRepo, animeRepo, userRepo)
	translationController := controllers.NewDescriptionTranslationController(translationService)

	r := chi.NewRouter()

	// Public
	r.Post("/anime/{animeID}", translationController.SubmitTranslation)
	r.Get("/anime/{animeID}", translationController.GetAnimeTranslation)
	r.Get("/user/{userID}", translationController.GetUserTranslations)

	// Requires moderator
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireRoleMiddleware(value.UserRoleModerator))
		r.Put("/{id}/accept", translationController.AcceptTranslation)
		r.Put("/{id}/reject", translationController.RejectTranslation)
		r.Get("/pending", translationController.GetPendingTranslations)
	})

	return r
}

func (a *Application) BootstrapUserModule() chi.Router {
	userRepo := repositories.NewUserRepository(a.Mongo)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := chi.NewRouter()

	// Public
	r.Get("/", userController.GetUsers)
	r.Get("/search", userController.SearchByUsername)
	r.Get("/{id}", userController.GetUserByID)

	// Requires authenticated
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))
		r.Put("/", userController.UpdateUserInfo)
	})

	// Requires moderator
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))
		r.Use(middlewares.RequireRoleMiddleware(value.UserRoleModerator))
		r.Put("/{id}/restrict", userController.RestrictAccount)
	})

	return r
}

func (a *Application) BootstrapAnimeModule() chi.Router {
	animeRepo := repositories.NewAnimeRepository()
	animeService := services.NewAnimeService(animeRepo)
	animeController := controllers.NewAnimeController(animeService)

	r := chi.NewRouter()
	r.Get("/{id}", animeController.GetAnimeByID)
	r.Get("/random", animeController.GetRandomAnime)
	r.Get("/search", animeController.SearchAnime)
	r.Get("/seasonal", animeController.GetAnimeThisSeason)
	r.Get("/studio/{id}", animeController.GetAnimeByStudioID)
	r.Get("/producer/{id}", animeController.GetAnimeByProducerID)
	r.Get("/licensor/{id}", animeController.GetAnimeByLicensorID)
	r.Get("/tags/{id}", animeController.GetAnimeByTagID)

	return r
}

func (a *Application) BootstrapFriendsModule() chi.Router {
	friendshipRepo := repositories.NewFriendshipRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipService := services.NewFriendshipService(userRepo, friendshipRepo)
	friendshipController := controllers.NewFriendshipController(friendshipService)

	r := chi.NewRouter()

	// Public
	r.Get("/{userID}", friendshipController.ListFriends)

	// Requires auth
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))
		r.Put("/send/{receiver}", friendshipController.SendFriendRequest)
		r.Put("/accept/{initiator}", friendshipController.AcceptFriendRequest)
		r.Put("/decline/{initiator}", friendshipController.DeclineFriendRequest)
		r.Put("/block/{receiver}", friendshipController.BlockUser)
		r.Get("/pending", friendshipController.ListPendingFriendRequests)
		r.Get("/check/{receiver}", friendshipController.FetchFriendshipStatus)
	})

	return r
}

func (a *Application) BootstrapAuthModule() chi.Router {

	jwtService := services.NewJWTService(a.JWTConfig)
	userService := services.NewUserService(repositories.NewUserRepository(a.Mongo))
	googleAuthController := controllers.NewGoogleAuthController(a.OAuth2Config, jwtService, userService)

	r := chi.NewRouter()

	// Rate limiter specific to authentication endpoints
	// This one must be stricter, as to avoid botted account creations
	authLimiter := &middlewares.IPRateLimiter{Rps: 0.5, Burst: 3}
	r.Use(authLimiter.Middleware)

	r.Get("/google/login", googleAuthController.Login)
	r.Get("/google/callback", googleAuthController.Callback)
	r.Get("/me", googleAuthController.WhoAmI)
	r.Get("/logout", googleAuthController.Logout)

	return r
}

func (a *Application) BootstrapReportsModule() chi.Router {
	reportRepo := repositories.NewUserReportRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	reportService := services.NewUserReportService(reportRepo, userRepo)
	reportController := controllers.NewUserReportController(reportService)

	r := chi.NewRouter()

	r.Post("/{userID}", reportController.SubmitReport)

	// Requires moderator
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireRoleMiddleware(value.UserRoleModerator))
		r.Get("/", reportController.GetReports)
		r.Get("/user/{userID}", reportController.GetReportsByTarget)
		r.Delete("/{id}", reportController.DeleteReport)
	})

	return r
}

func (a *Application) BootstrapPostModule() chi.Router {
	postRepo := repositories.NewPostRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipSvc := services.NewFriendshipService(userRepo, repositories.NewFriendshipRepository(a.Mongo))
	postService := services.NewPostService(postRepo, userRepo, friendshipSvc)
	postController := controllers.NewPostController(postService)

	r := chi.NewRouter()
	r.Get("/{post_id}", postController.GetPostById)
	r.Get("/{parent_id}/replies", postController.GetPostReplies)
	r.Post("/", postController.CreatePost)
	r.Post("/{post_id}/reply", postController.CreateReply)
	r.Delete("/{post_id}", postController.DeletePost)

	return r
}

func (a *Application) BootstrapRecommendationsModule() chi.Router {

	repo := repositories.NewRecommendationRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	service := services.NewRecommendationService(repo, userRepo)
	controller := controllers.NewRecommendationController(service)

	r := chi.NewRouter()
	r.Post("/{receiverID}/{animeID}", controller.Send)
	r.Get("/", controller.GetMine)
	r.Delete("/{animeID}", controller.Dismiss)

	return r
}

func (a *Application) BootstrapAnimeListModule() chi.Router {
	listRepo := repositories.NewAnimeListRepository(a.Mongo)
	animeRepo := repositories.NewAnimeRepository()

	listService := services.NewAnimeListService(listRepo, animeRepo)

	listController := controllers.NewAnimeListController(listService)

	r := chi.NewRouter()

	// Getting any user's list, completely public
	r.Get("/{userId}", listController.GetUserList)

	r.Group(func(r chi.Router) {
		// Requires auth because the user is editing their own list
		r.Use(middlewares.JWTMiddleware(a.JWTConfig))

		r.Post("/{userId}/{animeId}", listController.AddAnime)
		r.Patch("/{userId}/progress/{animeId}", listController.UpdateProgress)
		r.Patch("/{userId}/status/{animeId}", listController.UpdateStatus)
		r.Patch("/{userId}/notes/{animeId}", listController.UpdateNotes)
		r.Patch("/{userId}/rating/{animeId}", listController.UpdateRating)
		r.Delete("/{userId}/{animeId}", listController.RemoveAnimeFromList)
	})

	return r
}
