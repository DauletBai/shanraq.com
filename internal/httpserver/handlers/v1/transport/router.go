package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
)

// Router configures routes for transportation logistics partners.
func Router(cfg config.Config, logger zerolog.Logger) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("filter_country", r.URL.Query().Get("country")).Msg("transport_company_list")
		writeJSON(w, http.StatusOK, map[string]any{
			"data":  []any{},
			"meta":  map[string]any{"message": "transport companies listing to be implemented"},
			"flags": map[string]bool{"enabled": cfg.Features.EnableTransportCompanies},
		})
	})

	r.Post("/", notYetImplemented)
	r.Get("/{id}", notYetImplemented)
	r.Put("/{id}", notYetImplemented)
	r.Delete("/{id}", notYetImplemented)

	return r
}

func notYetImplemented(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]any{
		"error":   "not_implemented",
		"message": "transport company endpoint planned for upcoming iteration",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
