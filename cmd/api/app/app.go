package app

import (
	"log"

	"github.com/afuradanime/backend/cmd/api/app/database"
	"github.com/afuradanime/backend/config"
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type Application struct {
	Config       *config.Config
	OAuth2Config *oauth2.Config
	JWTConfig    *config.JWTConfig
	Mongo        *mongo.Database // The mongo database handle
}

func New() *Application {

	env := utils.GetApplicationEnvironment()
	log.Println("Starting application in " + env + " environment...")

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

	if Config.ShouldBootstrap || env == "test" /* Always bootstrap on test */ {
		log.Println("Bootstrapping database...")
		Bootstrap(app.Mongo)
	}

	return app
}

func (app *Application) Run() {

	s := fuego.NewServer(
		fuego.WithAddr("localhost:"+app.Config.Port),
		fuego.WithGlobalMiddlewares(middlewares.CORSMiddleware),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				SwaggerURL:   "/swagger",
				SpecURL:      "/swagger/openapi.json",
				JSONFilePath: "_openapi/openapi.json",
				Info: &openapi3.Info{
					Title:       "Afuradanime API",
					Description: "The openapi docs for the Afuradanime API/Backend",
					Version:     "0.2.0",
				},
			}),
		),
	)

	app.InitRoutes(s)

	log.Println("Server started on port " + app.Config.Port)
	err := s.Run()
	if err != nil {
		log.Fatal("Actually, the server failed to start: ", err)
	}
}
