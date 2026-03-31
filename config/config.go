package config

import (
	"log"
	"os"

	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	FrontendURL           string
	ShouldBootstrap       bool
	AnimeDatabasePath     string
	MongoConnectionString string
	MongoUsername         string
	MongoPassword         string
	MongoDatabase         string
	EncryptionKey         string
}

func Load() *Config {
	env := utils.GetApplicationEnvironment()

	// Load the appropriate .env file based on the environment
	envFile := ".env"
	if env == "test" {
		envFile = ".env.test"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Println("Warning: no .env file found at", envFile)
		panic("No .env file")
	}

	return &Config{
		Port:                  os.Getenv("PORT"),
		FrontendURL:           os.Getenv("FRONTEND_URL"),
		ShouldBootstrap:       os.Getenv("SHOULD_BOOTSTRAP") == "true",
		AnimeDatabasePath:     os.Getenv("ANIME_DATABASE_PATH"),
		MongoConnectionString: os.Getenv("MONGO_CONNECTION_STRING"),
		MongoUsername:         os.Getenv("MONGO_USERNAME"),
		MongoPassword:         os.Getenv("MONGO_PASSWORD"),
		MongoDatabase:         os.Getenv("MONGO_DATABASE"),
		EncryptionKey:         os.Getenv("ENCRYPTION_KET"),
	}
}
