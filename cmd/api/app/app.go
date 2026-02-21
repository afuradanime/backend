package app

import (
	"log"
	"net/http"

	"github.com/afuradanime/backend/cmd/api/app/database"
	"github.com/afuradanime/backend/config"
	"github.com/afuradanime/backend/internal/core/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	Config       *config.Config
	OAuth2Config *oauth2.Config
	JWTConfig    *config.JWTConfig
	Mongo        *mongo.Database // The mongo database handle
}

func New() *Application {

	Config := config.Load()
	OAuth2 := config.LoadOauth2()
	JWTConfig := config.LoadJWTConfig()
	log.Println("Config loaded successfully!")

	// Setup Sqlite connection
	database.InitSQLite(*Config)

	// Start MongoDB connection
	mongoClient, err := database.InitMongoDB(*Config)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB: ", err)
	}

	// Start encryption service
	utils.InitEncryption(Config.EncryptionKey)

	app := &Application{
		Mongo:        mongoClient.Database(Config.MongoDatabase),
		Config:       Config,
		OAuth2Config: OAuth2,
		JWTConfig:    JWTConfig,
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
	app.InitRoutes(r)

	log.Println("Server started on port " + app.Config.Port)
	err := http.ListenAndServe(":"+app.Config.Port, r)
	if err != nil {
		log.Fatal("Actually, the server failed to start: ", err)
	}
}
