package config

import (
	"fmt"
	"log"
	"os"

	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/joho/godotenv"
)

type JWTConfig struct {
	Secret          string
	ExpirationHours int64
	Issuer          string
}

func LoadJWTConfig() *JWTConfig {

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

	secret := os.Getenv("JWT_SECRET")
	expirationHours := JWTExpToInt64(os.Getenv("JWT_EXPIRATION_HOURS"))
	issuer := os.Getenv("JWT_ISSUER")

	return &JWTConfig{
		Secret:          secret,
		ExpirationHours: expirationHours,
		Issuer:          issuer,
	}

	// return &JWTConfig{
	// 	Secret:          os.Getenv("JWT_SECRET"),
	// 	ExpirationHours: JWTExpToInt64(os.Getenv("JWT_EXPIRATION_HOURS")),
	// 	Issuer:          os.Getenv("JWT_ISSUER"),
	// }
}

func JWTExpToInt64(exp string) int64 {
	var expInt int64
	_, err := fmt.Sscanf(exp, "%d", &expInt)
	if err != nil {
		log.Fatalf("Error converting JWT_EXPIRATION_HOURS to int64: %v", err)
	}
	return expInt
}
