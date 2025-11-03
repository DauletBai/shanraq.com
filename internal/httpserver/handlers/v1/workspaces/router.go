package workspaces

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/auth/session"
	"shanraq.com/internal/config"
	workspaceservice "shanraq.com/internal/services/workspace"
)

// Router exposes workspace APIs for authenticated users.
func Router(cfg config.Config, logger zerolog.Logger, svc workspaceservice.Service) chi.Router {
	_ = cfg
	r := chi.NewRouter()

	r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		identity, ok := session.IdentityFromContext(r.Context())
		if !ok {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
			return
		}
		workspace, err := svc.GetOrCreate(r.Context(), identity)
		if err != nil {
			logger.Error().Err(err).Str("user", identity.Subject).Msg("workspace_get")
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "workspace_error"})
			return
		}
		respondJSON(w, http.StatusOK, workspace)
	})

	r.Post("/me/plans", func(w http.ResponseWriter, r *http.Request) {
		identity, ok := session.IdentityFromContext(r.Context())
		if !ok {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
			return
		}
		var plan workspaceservice.BusinessPlan
		if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_payload"})
			return
		}
		workspace, err := svc.AddPlan(r.Context(), identity, plan)
		if err != nil {
			logger.Error().Err(err).Str("user", identity.Subject).Msg("workspace_add_plan")
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "workspace_error"})
			return
		}
		respondJSON(w, http.StatusCreated, workspace)
	})

	return r
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
