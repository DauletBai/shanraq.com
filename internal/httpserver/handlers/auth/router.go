package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/auth"
	"shanraq.com/internal/auth/session"
	"shanraq.com/internal/config"
)

// Router exposes auth related endpoints (login, callback, logout) for multiple providers.
func Router(cfg config.Config, logger zerolog.Logger, registry *auth.ProviderRegistry, sessions *session.Manager) chi.Router {
	if registry == nil {
		registry = auth.NewRegistry()
	}

	r := chi.NewRouter()

	r.Get("/providers", func(w http.ResponseWriter, _ *http.Request) {
		respondJSON(w, http.StatusOK, map[string]any{
			"providers": registry.List(),
		})
	})

	r.Route("/{provider}", func(pr chi.Router) {
		pr.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			providerName := chi.URLParam(r, "provider")
			provider, err := registry.Get(providerName)
			if err != nil {
				logger.Warn().Str("provider", providerName).Err(err).Msg("auth_provider_not_available")
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "provider_not_configured",
				})
				return
			}

			state := r.URL.Query().Get("state")
			if state == "" {
				state = generateState()
			}

			authURL, err := provider.AuthCodeURL(state)
			if err != nil {
				logger.Warn().Str("provider", providerName).Err(err).Msg("auth_login_not_configured")
				respondJSON(w, http.StatusNotImplemented, map[string]any{
					"error":   "auth_not_configured",
					"message": "Authentication provider is not configured yet.",
				})
				return
			}

			http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
		})

		pr.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
			providerName := chi.URLParam(r, "provider")
			provider, err := registry.Get(providerName)
			if err != nil {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "provider_not_configured",
				})
				return
			}

			code := r.URL.Query().Get("code")
			state := r.URL.Query().Get("state")
			if code == "" {
				respondJSON(w, http.StatusBadRequest, map[string]string{"error": "missing_code"})
				return
			}

			identity, err := provider.Exchange(r.Context(), code)
			if err != nil {
				logger.Error().Str("provider", providerName).Err(err).Msg("auth_exchange_failed")
				respondJSON(w, http.StatusBadGateway, map[string]string{"error": "exchange_failed"})
				return
			}

			if sessions != nil {
				if _, err := sessions.Create(w, identity); err != nil {
					logger.Error().Err(err).Msg("create_session_failed")
					respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "session_error"})
					return
				}
			}

			respondJSON(w, http.StatusOK, map[string]any{
				"state":    state,
				"identity": identity,
			})
		})
	})

	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		returnURL := r.URL.Query().Get("return")
		if returnURL == "" {
			returnURL = cfg.HTTP.PublicBaseURL
		}
		if _, err := url.Parse(returnURL); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_return_url"})
			return
		}
		if sessions != nil {
			sessions.Destroy(w, r)
		}
		respondJSON(w, http.StatusOK, map[string]string{
			"message":    "logged_out",
			"return_url": returnURL,
		})
	})

	r.Get("/session", func(w http.ResponseWriter, r *http.Request) {
		if sessions == nil {
			respondJSON(w, http.StatusOK, map[string]any{"authenticated": false})
			return
		}
		identity, ok := sessions.Identity(r)
		respondJSON(w, http.StatusOK, map[string]any{
			"authenticated": ok,
			"identity":      identity,
		})
	})

	return r
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func generateState() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
