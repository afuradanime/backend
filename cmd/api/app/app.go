package app

import (
	"log"
	"net/http"

	"github.com/afuradanime/backend/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New() {
	Config := config.Load()
	log.Println("Config loaded successfully!")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{Config.FrontendURL, "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	InitRoutes(r)

	log.Println("Server started on port " + Config.Port)
	err := http.ListenAndServe(":"+Config.Port, r)
	if err != nil {
		log.Fatal("Actually, the server failed to start: ", err)
	}
}
