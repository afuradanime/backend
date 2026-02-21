package config

import (
	"log"
	"os"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
