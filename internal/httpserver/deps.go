package httpserver

import (
	"github.com/rs/zerolog"

	"shanraq.com/internal/auth"
	"shanraq.com/internal/config"
	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
	"shanraq.com/internal/web"
)

// Deps aggregates dependencies required to build the HTTP router.
type Deps struct {
	Logger           zerolog.Logger
	Config           config.Config
	Renderer         *web.Renderer
	TransportService transportservice.Service
	AgencyService    agencyservice.Service
	ListingService   listingservice.Service
	AuthRegistry     *auth.ProviderRegistry
}
