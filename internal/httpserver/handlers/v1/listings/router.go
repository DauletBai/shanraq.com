package listings

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	listingservice "shanraq.com/internal/services/listing"
)

// Router exposes property listing read endpoints.
func Router(cfg config.Config, logger zerolog.Logger, svc listingservice.Service) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		listings, err := svc.List(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("list_listings_failed")
			respondError(w, http.StatusInternalServerError, "list_failed")
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"data": listings,
			"meta": map[string]any{
				"count": len(listings),
			},
		})
	})

	r.Get("/featured", func(w http.ResponseWriter, r *http.Request) {
		listings, err := svc.Featured(r.Context(), 6)
		if err != nil {
			logger.Error().Err(err).Msg("featured_listings_failed")
			respondError(w, http.StatusInternalServerError, "list_failed")
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"data": listings,
		})
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id")
			return
		}
		listing, err := svc.Get(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not_found")
			return
		}
		respondJSON(w, http.StatusOK, listing)
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
