package utils

import "github.com/golang-jwt/jwt/v5"

func GetParsedJWTClaims(tokenString string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(""), nil
	})

	if err != nil || !parsed.Valid {
		return nil, err
	}

	return parsed, nil
}
