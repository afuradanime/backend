package middlewares

import (
	"context"
	"net/http"

	"github.com/afuradanime/backend/config"
	"github.com/afuradanime/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey    contextKey = "userID"
	UserRolesKey contextKey = "userRoles"
)

func JWTMiddleware(cfg *config.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("jwt")
			if err != nil {
				http.Error(w, "Unauthorized, no token provided", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			if tokenString == "" {
				http.Error(w, "Unauthorized, empty token", http.StatusUnauthorized)
				return
			}

			parsed, err := utils.GetParsedJWTClaims(tokenString, cfg.Secret)
			if err != nil || !parsed.Valid {
				http.Error(w, "Unauthorized, invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := parsed.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized, invalid claims", http.StatusUnauthorized)
				return
			}

			idValue, ok := claims["id"]
			if !ok {
				http.Error(w, "Unauthorized, missing user ID in token", http.StatusUnauthorized)
				return
			}

			userID, ok := idValue.(float64)
			if !ok {
				http.Error(w, "Unauthorized, invalid user ID in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
			roles, ok := claims["role"]
			if ok {
				ctx = context.WithValue(ctx, UserRolesKey, roles)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
