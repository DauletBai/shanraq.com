package middleware

import (
	"net/http"

	"shanraq.com/internal/models"
)

func Authorize(alloweRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
			user, ok := r.Context().Value("user").(*models.User)
			if !ok {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}
			hasRole := false
			for _, role := range alloweRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}