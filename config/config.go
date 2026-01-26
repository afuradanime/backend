package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	FrontendURL string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}
}
