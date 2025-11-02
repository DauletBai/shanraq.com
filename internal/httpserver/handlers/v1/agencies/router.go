package agencies

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	agencyservice "shanraq.com/internal/services/agency"
)

// Router exposes agency and realtor read endpoints.
func Router(cfg config.Config, logger zerolog.Logger, svc agencyservice.Service) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		agencies, err := svc.ListAgencies(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("list_agencies_failed")
			respondError(w, http.StatusInternalServerError, "list_failed")
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"data": agencies,
			"meta": map[string]any{
				"count":          len(agencies),
				"country_filter": cfg.Geo.DataProvider,
			},
		})
	})

	r.Get("/realtors", func(w http.ResponseWriter, r *http.Request) {
		realtors, err := svc.ListRealtors(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("list_realtors_failed")
			respondError(w, http.StatusInternalServerError, "list_failed")
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"data": realtors,
			"meta": map[string]any{
				"count": len(realtors),
			},
		})
	})

	r.Get("/realtors/featured", func(w http.ResponseWriter, r *http.Request) {
		realtors, err := svc.FeaturedRealtors(r.Context(), 4)
		if err != nil {
			logger.Error().Err(err).Msg("featured_realtors_failed")
			respondError(w, http.StatusInternalServerError, "list_failed")
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"data": realtors,
		})
	})

	return r
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, code string) {
	respondJSON(w, status, map[string]string{"error": code})
}
