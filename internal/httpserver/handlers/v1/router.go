package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	"shanraq.com/internal/httpserver/handlers/v1/agencies"
	"shanraq.com/internal/httpserver/handlers/v1/listings"
	"shanraq.com/internal/httpserver/handlers/v1/transport"
	"shanraq.com/internal/httpserver/handlers/v1/workspaces"
	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
	workspaceservice "shanraq.com/internal/services/workspace"
)

// Router wires REST API routes under /api/v1.
func Router(cfg config.Config, logger zerolog.Logger, transportSvc transportservice.Service, agencySvc agencyservice.Service, listingSvc listingservice.Service, workspaceSvc workspaceservice.Service) chi.Router {
	r := chi.NewRouter()

	r.Mount("/transport-companies", transport.Router(cfg, logger, transportSvc))
	r.Mount("/agencies", agencies.Router(cfg, logger, agencySvc))
	r.Mount("/listings", listings.Router(cfg, logger, listingSvc))
	r.Mount("/workspaces", workspaces.Router(cfg, logger, workspaceSvc))

	return r
}
