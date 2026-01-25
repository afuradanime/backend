package app

import (
	"log"

	"github.com/afuradanime/backend/internal/adapters/controllers"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r *chi.Mux) {
	r.Mount("/users", BootstrapUserModule())
	log.Println("[Routing] User routes initialized...")

	//...

	log.Println("[Routing] All routes initialized successfully!")
}

func BootstrapUserModule() chi.Router {
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := chi.NewRouter()
	r.Get("/{id}", userController.GetUserByID)

	return r
}
