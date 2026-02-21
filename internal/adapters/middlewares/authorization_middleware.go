package middlewares

import (
	"net/http"

	"github.com/afuradanime/backend/internal/core/domain/value"
)

func RequireRoleMiddleware(role value.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !IsLoggedUserOfRole(r, role) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
