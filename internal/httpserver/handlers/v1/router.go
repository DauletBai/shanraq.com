package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	"shanraq.com/internal/httpserver/handlers/v1/transport"
)

// Router wires REST API routes under /api/v1.
func Router(cfg config.Config, logger zerolog.Logger) chi.Router {
	r := chi.NewRouter()

	r.Mount("/transport-companies", transport.Router(cfg, logger))

	return r
}
