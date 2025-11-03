package session

import (
	"net/http"
)

// Middleware resolves sessions for incoming requests and injects identities into the context.
func Middleware(manager *Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if manager != nil {
				if identity, ok := manager.Identity(r); ok {
					r = r.WithContext(WithIdentity(r.Context(), identity))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
