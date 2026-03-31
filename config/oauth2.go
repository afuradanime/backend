package config

import (
	"log"
	"os"

	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func LoadOauth2() *oauth2.Config {

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

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",   // Required for email access
			"https://www.googleapis.com/auth/userinfo.profile", // Required for profile access (name, picture, etc...
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}
