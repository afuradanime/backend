package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	FrontendURL           string
	AnimeDatabasePath     string
	MongoConnectionString string
	MongoUsername         string
	MongoPassword         string
	ShouldBootstrap       bool
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:                  os.Getenv("PORT"),
		FrontendURL:           os.Getenv("FRONTEND_URL"),
		AnimeDatabasePath:     os.Getenv("ANIME_DATABASE_PATH"),
		MongoConnectionString: os.Getenv("MONGO_CONNECTION_STRING"),
		MongoUsername:         os.Getenv("MONGO_USERNAME"),
		MongoPassword:         os.Getenv("MONGO_PASSWORD"),
		ShouldBootstrap:       os.Getenv("SHOULD_BOOTSTRAP") == "true",
	}
}
