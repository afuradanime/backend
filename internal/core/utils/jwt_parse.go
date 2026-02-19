package utils

import (
	"encoding/base64"

	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/golang-jwt/jwt/v5"
)

func GetParsedJWTClaims(tokenString string, secret string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil || !parsed.Valid {
		return nil, err
	}

	return parsed, nil
}

// We store the roles as strings, so we get the claims back as "Aw=="
// instead of [3] for some reason, I don't really know why, but we can just convert it easy
// so it's probably fine
func DecodeRoleList(encodedRoles string) []value.UserRole {
	decoded, err := base64.StdEncoding.DecodeString(encodedRoles)
	if err != nil {
		return []value.UserRole{}
	}

	roles := make([]value.UserRole, len(decoded))
	for i, b := range decoded {
		roles[i] = value.UserRole(b)
	}
	return roles
}
