package app

import (
	"log"
	"net/http"

	"github.com/afuradanime/backend/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() {
	Config := config.Load()
	log.Println("Config loaded successfully!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	InitRoutes(r)

	log.Println("Server started on port " + Config.Port)
	err := http.ListenAndServe(":"+Config.Port, r)
	if err != nil {
		log.Fatal("Actually, the server failed to start: ", err)
	}
}
