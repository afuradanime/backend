package config

import (
	"fmt"
	"log"
	"os"
)

type JWTConfig struct {
	Secret          string
	ExpirationHours int64
	Issuer          string
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:          os.Getenv("JWT_SECRET"),
		ExpirationHours: JWTExpToInt64(os.Getenv("JWT_EXPIRATION_HOURS")),
		Issuer:          os.Getenv("JWT_ISSUER"),
	}
}

func JWTExpToInt64(exp string) int64 {
	var expInt int64
	_, err := fmt.Sscanf(exp, "%d", &expInt)
	if err != nil {
		log.Fatalf("Error converting JWT_EXPIRATION_HOURS to int64: %v", err)
	}
	return expInt
}
