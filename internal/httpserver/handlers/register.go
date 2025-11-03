package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/auth"
	"shanraq.com/internal/auth/session"
	"shanraq.com/internal/config"
	authhandler "shanraq.com/internal/httpserver/handlers/auth"
	"shanraq.com/internal/httpserver/handlers/public"
	"shanraq.com/internal/httpserver/handlers/v1"
	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
	workspaceservice "shanraq.com/internal/services/workspace"
	"shanraq.com/internal/web"
)

// RegisterRoutes configures the HTTP routes.
func RegisterRoutes(
	r chi.Router,
	cfg config.Config,
	logger zerolog.Logger,
	renderer *web.Renderer,
	transportSvc transportservice.Service,
	agencySvc agencyservice.Service,
	listingSvc listingservice.Service,
	authRegistry *auth.ProviderRegistry,
	sessionManager *session.Manager,
	workspaceSvc workspaceservice.Service,
) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Mount("/", public.Router(cfg, logger, renderer, listingSvc, agencySvc, transportSvc))
	r.Mount("/api/v1", v1.Router(cfg, logger, transportSvc, agencySvc, listingSvc, workspaceSvc))
	r.Mount("/auth", authhandler.Router(cfg, logger, authRegistry, sessionManager))
}
