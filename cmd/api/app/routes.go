package app

import (
	"github.com/afuradanime/backend/internal/adapters/controllers"
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
)

func (a *Application) InitRoutes(s *fuego.Server) {
	// Global rate limiter
	globalLimiter := &middlewares.IPRateLimiter{Rps: 10, Burst: 30}

	// Fuego uses package level Use function
	fuego.Use(s,
		middleware.Logger,
		middleware.Recoverer,
		// middlewares.CORSMiddleware,
		globalLimiter.Middleware,
	)

	// Register Modules
	a.RegisterAuthModule(s)
	a.RegisterAnimeModule(s)
	a.RegisterUserModule(s)
	a.RegisterFriendsModule(s)
	a.RegisterTranslationsModule(s)
	a.RegisterAnimeListModule(s)

	// Group for globally protected routes
	protected := fuego.Group(s, "/")
	fuego.Use(protected, middlewares.JWTMiddleware(a.JWTConfig))

	a.RegisterReportsModule(protected)
	a.RegisterPostModule(protected)
	a.RegisterRecommendationsModule(protected)
}

func (a *Application) RegisterTranslationsModule(s *fuego.Server) {
	translationRepo := repositories.NewDescriptionTranslationRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	animeRepo := repositories.NewAnimeRepository()
	translationService := services.NewDescriptionTranslationService(translationRepo, animeRepo, userRepo)
	translationController := controllers.NewDescriptionTranslationController(translationService)

	g := fuego.Group(s, "/translations")

	// Public
	fuego.Post(g, "/anime/{animeID}", translationController.SubmitTranslation)
	fuego.Get(g, "/anime/{animeID}", translationController.GetAnimeTranslation)
	fuego.Get(g, "/user/{userID}", translationController.GetUserTranslations)

	// Moderator only
	modGroup := fuego.Group(g, "/")
	fuego.Use(modGroup, middlewares.RequireRoleMiddleware(value.UserRoleModerator))
	fuego.Put(modGroup, "/{id}/accept", translationController.AcceptTranslation)
	fuego.Put(modGroup, "/{id}/reject", translationController.RejectTranslation)
	fuego.Get(modGroup, "/pending", translationController.GetPendingTranslations)
}

func (a *Application) RegisterUserModule(s *fuego.Server) {
	userRepo := repositories.NewUserRepository(a.Mongo)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	g := fuego.Group(s, "/users")

	fuego.Get(g, "/", userController.GetUsers)
	fuego.Get(g, "/search", userController.SearchByUsername)
	fuego.Get(g, "/{id}", userController.GetUserByID)

	// Authenticated
	authGroup := fuego.Group(g, "/")
	fuego.Use(authGroup, middlewares.JWTMiddleware(a.JWTConfig))
	fuego.Put(authGroup, "/", userController.UpdateUserInfo)

	// Moderator
	modGroup := fuego.Group(authGroup, "/")
	fuego.Use(modGroup, middlewares.RequireRoleMiddleware(value.UserRoleModerator))
	fuego.Put(modGroup, "/{id}/restrict", userController.RestrictAccount)
}

func (a *Application) RegisterAnimeModule(s *fuego.Server) {
	animeRepo := repositories.NewAnimeRepository()
	animeService := services.NewAnimeService(animeRepo)
	animeController := controllers.NewAnimeController(animeService)

	g := fuego.Group(s, "/anime")
	fuego.Get(g, "/{id}", animeController.GetAnimeByID)
	fuego.Get(g, "/random", animeController.GetRandomAnime)
	fuego.Get(g, "/search", animeController.SearchAnime)
	fuego.Get(g, "/seasonal", animeController.GetAnimeThisSeason)
	fuego.Get(g, "/studio/{id}", animeController.GetAnimeByStudioID)
	fuego.Get(g, "/producer/{id}", animeController.GetAnimeByProducerID)
	fuego.Get(g, "/licensor/{id}", animeController.GetAnimeByLicensorID)
	fuego.Get(g, "/tags/{id}", animeController.GetAnimeByTagID)
}

func (a *Application) RegisterFriendsModule(s *fuego.Server) {
	friendshipRepo := repositories.NewFriendshipRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipService := services.NewFriendshipService(userRepo, friendshipRepo)
	friendshipController := controllers.NewFriendshipController(friendshipService)

	g := fuego.Group(s, "/friends")
	fuego.Get(g, "/{userID}", friendshipController.ListFriends)

	authGroup := fuego.Group(g, "/")
	fuego.Use(authGroup, middlewares.JWTMiddleware(a.JWTConfig))
	fuego.Put(authGroup, "/send/{receiver}", friendshipController.SendFriendRequest)
	fuego.Put(authGroup, "/accept/{initiator}", friendshipController.AcceptFriendRequest)
	fuego.Put(authGroup, "/decline/{initiator}", friendshipController.DeclineFriendRequest)
	fuego.Put(authGroup, "/block/{receiver}", friendshipController.BlockUser)
	fuego.Get(authGroup, "/pending", friendshipController.ListPendingFriendRequests)
	fuego.Get(authGroup, "/check/{receiver}", friendshipController.FetchFriendshipStatus)
}

func (a *Application) RegisterAuthModule(s *fuego.Server) {
	jwtService := services.NewJWTService(a.JWTConfig)
	userService := services.NewUserService(repositories.NewUserRepository(a.Mongo))
	googleAuthController := controllers.NewGoogleAuthController(a.Config, a.OAuth2Config, jwtService, userService)

	g := fuego.Group(s, "/auth")
	authLimiter := &middlewares.IPRateLimiter{Rps: 0.5, Burst: 3}
	fuego.Use(g, authLimiter.Middleware)

	fuego.Get(g, "/google/login", googleAuthController.Login)
	fuego.Get(g, "/google/callback", googleAuthController.Callback)
	fuego.Get(g, "/me", googleAuthController.WhoAmI)
	fuego.Get(g, "/logout", googleAuthController.Logout)
}

func (a *Application) RegisterReportsModule(s *fuego.Server) {
	reportRepo := repositories.NewUserReportRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	reportService := services.NewUserReportService(reportRepo, userRepo)
	reportController := controllers.NewUserReportController(reportService)

	g := fuego.Group(s, "/reports")
	fuego.Post(g, "/{userID}", reportController.SubmitReport)

	modGroup := fuego.Group(g, "/")
	fuego.Use(modGroup, middlewares.RequireRoleMiddleware(value.UserRoleModerator))
	fuego.Get(modGroup, "/", reportController.GetReports)
	fuego.Get(modGroup, "/user/{userID}", reportController.GetReportsByTarget)
	fuego.Delete(modGroup, "/{id}", reportController.DeleteReport)
}

func (a *Application) RegisterPostModule(s *fuego.Server) {
	postRepo := repositories.NewPostRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	friendshipSvc := services.NewFriendshipService(userRepo, repositories.NewFriendshipRepository(a.Mongo))
	postService := services.NewPostService(postRepo, userRepo, friendshipSvc)
	postController := controllers.NewPostController(postService)

	g := fuego.Group(s, "/posts")
	fuego.Get(g, "/{post_id}", postController.GetPostById)
	fuego.Get(g, "/{parent_id}/replies", postController.GetPostReplies)
	fuego.Post(g, "/", postController.CreatePost)
	fuego.Post(g, "/{post_id}/reply", postController.CreateReply)
	fuego.Delete(g, "/{post_id}", postController.DeletePost)
}

func (a *Application) RegisterRecommendationsModule(s *fuego.Server) {
	repo := repositories.NewRecommendationRepository(a.Mongo)
	userRepo := repositories.NewUserRepository(a.Mongo)
	service := services.NewRecommendationService(repo, userRepo)
	controller := controllers.NewRecommendationController(service)

	g := fuego.Group(s, "/recommendations")
	fuego.Post(g, "/{receiverID}/{animeID}", controller.Send)
	fuego.Get(g, "/", controller.GetMine)
	fuego.Delete(g, "/{animeID}", controller.Dismiss)
}

func (a *Application) RegisterAnimeListModule(s *fuego.Server) {
	listRepo := repositories.NewAnimeListRepository(a.Mongo)
	animeRepo := repositories.NewAnimeRepository()
	listService := services.NewAnimeListService(listRepo, animeRepo)
	listController := controllers.NewAnimeListController(listService)

	g := fuego.Group(s, "/animelist")
	fuego.Get(g, "/{userId}", listController.GetUserList)

	authGroup := fuego.Group(g, "/")
	fuego.Use(authGroup, middlewares.JWTMiddleware(a.JWTConfig))
	fuego.Post(authGroup, "/{userId}/{animeId}", listController.AddAnime)
	fuego.Patch(authGroup, "/{userId}/progress/{animeId}", listController.UpdateProgress)
	fuego.Patch(authGroup, "/{userId}/status/{animeId}", listController.UpdateStatus)
	fuego.Patch(authGroup, "/{userId}/notes/{animeId}", listController.UpdateNotes)
	fuego.Patch(authGroup, "/{userId}/rating/{animeId}", listController.UpdateRating)
	fuego.Delete(authGroup, "/{userId}/{animeId}", listController.RemoveAnimeFromList)
}
