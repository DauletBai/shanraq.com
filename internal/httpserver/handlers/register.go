package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	"shanraq.com/internal/httpserver/handlers/public"
	"shanraq.com/internal/httpserver/handlers/v1"
	"shanraq.com/internal/web"
)

// RegisterRoutes configures the HTTP routes.
func RegisterRoutes(r chi.Router, cfg config.Config, logger zerolog.Logger, renderer *web.Renderer) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Mount("/", public.Router(cfg, logger, renderer))
	r.Mount("/api/v1", v1.Router(cfg, logger))
}
