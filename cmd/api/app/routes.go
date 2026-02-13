package app

import (
	"log"

	"github.com/afuradanime/backend/internal/adapters/controllers"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-chi/chi/v5"
)

func (a *Application) InitRoutes(r *chi.Mux) {
	r.Mount("/users", a.BootstrapUserModule())
	log.Println("[Routing] User routes initialized...")

	r.Mount("/anime", a.BootstrapAnimeModule())
	log.Println("[Routing] Anime routes initialized...")

	log.Println("[Routing] All routes initialized successfully!")
}

func (a *Application) BootstrapUserModule() chi.Router {
	userRepo := repositories.NewUserRepository(a.Mongo)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := chi.NewRouter()
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

	return r
}
