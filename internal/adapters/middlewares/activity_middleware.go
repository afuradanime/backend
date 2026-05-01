package middlewares

import (
	"net/http"

	"github.com/afuradanime/backend/config"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/golang-jwt/jwt/v5"
)

func ActivityMiddleware(cfg *config.JWTConfig, tracker *domain.ActivityTracker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err == nil && cookie.Value != "" {
				parsed, err := utils.GetParsedJWTClaims(cookie.Value, cfg.Secret)
				if err == nil && parsed.Valid {
					if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
						if idValue, ok := claims["id"]; ok {
							if userID, ok := idValue.(float64); ok {
								tracker.RecordActivity(int(userID), value.Online)
							}
						}
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
