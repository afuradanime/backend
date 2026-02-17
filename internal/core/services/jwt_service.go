package services

import (
	"time"

	"github.com/afuradanime/backend/config"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Super secret, architecture breaking, JWT generation service!
type JWTService struct {
	conf *config.JWTConfig
}

func NewJWTService(config *config.JWTConfig) *JWTService {
	return &JWTService{conf: config}
}

func (s *JWTService) GenerateJWT(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":   user.ID,
		"iss":  s.conf.Issuer,
		"exp":  time.Now().Add(time.Hour * time.Duration(s.conf.ExpirationHours)).Unix(),
		"role": user.Roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.conf.Secret))
}

func (s *JWTService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	parsed, err := utils.GetParsedJWTClaims(tokenString, s.conf.Secret)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
