package public

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
	"shanraq.com/internal/web"
)

// Router exposes the public-facing endpoints.
func Router(
	cfg config.Config,
	logger zerolog.Logger,
	renderer *web.Renderer,
	listingSvc listingservice.Service,
	agencySvc agencyservice.Service,
	transportSvc transportservice.Service,
) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("path", r.URL.Path).Msg("public_page")
		w.Header().Set("X-App-Name", cfg.App.Name)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if renderer != nil {
			data := &web.HomePageData{}
			data.BrandName = cfg.App.Name
			data.PageTitle = "Global Real Estate Platform · "
			data.Description = "Discover, list, and manage properties and logistics partners across the world with Shanraq."
			data.PageID = "home"

			if listingSvc != nil {
				featuredListings, err := listingSvc.Featured(r.Context(), 6)
				if err != nil {
					logger.Warn().Err(err).Msg("fetch_featured_listings")
				} else {
					data.FeaturedListings = web.MapListings(featuredListings)
				}
			}
			if agencySvc != nil {
				agencies, err := agencySvc.Featured(r.Context(), 3)
				if err != nil {
					logger.Warn().Err(err).Msg("fetch_featured_agencies")
				} else {
					data.FeaturedAgencies = web.MapAgencies(agencies)
				}
				if realtors, err := agencySvc.FeaturedRealtors(r.Context(), 4); err == nil {
					data.FeaturedRealtors = web.MapRealtors(realtors)
				}
			}
			if transportSvc != nil {
				companies, _, err := transportSvc.List(r.Context(), transportservice.ListFilter{ActiveOnly: true, Limit: 4})
				if err != nil {
					logger.Warn().Err(err).Msg("fetch_featured_transport")
				} else {
					data.FeaturedTransport = web.MapTransportCompanies(companies)
				}
			}

			if err := renderer.RenderHome(w, data); err != nil {
				logger.Error().Err(err).Msg("render_home")
				http.Error(w, "unable to render", http.StatusInternalServerError)
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Welcome to Shanraq Real Estate"))
	})

	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/", http.StatusTemporaryRedirect)
	})

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))
	r.Handle("/static/*", staticHandler)

	dashboardHandler := http.StripPrefix("/dashboard/", http.FileServer(http.Dir("web/dashboard")))
	r.Handle("/dashboard/*", dashboardHandler)

	return r
}
