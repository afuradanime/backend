package app

import (
	"log"
	"net/http"

	"github.com/afuradanime/backend/cmd/api/app/database"
	"github.com/afuradanime/backend/config"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Application struct {
	Config *config.Config
	Mongo  *mongo.Database // The mongo database handle
}

func New() *Application {

	Config := config.Load()
	log.Println("Config loaded successfully!")

	// Setup Sqlite connection
	database.InitSQLite(*Config)

	// Start MongoDB connection
	mongoClient, err := database.InitMongoDB(*Config)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB: ", err)
	}

	app := &Application{
		Mongo:  mongoClient.Database("afuradanime"),
		Config: Config,
	}

	if Config.ShouldBootstrap {
		log.Println("Bootstrapping database...")
		app.Bootstrap()
	}

	return app
}

func (app *Application) Run() {
	// Setup HTTP server
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer, // useful middleware to recover from panics and return a 500 error
	)

	// CORS setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.Config.FrontendURL, "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.InitRoutes(r)

	log.Println("Server started on port " + app.Config.Port)
	err := http.ListenAndServe(":"+app.Config.Port, r)
	if err != nil {
		log.Fatal("Actually, the server failed to start: ", err)
	}
}
