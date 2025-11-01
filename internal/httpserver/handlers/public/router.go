package public

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	"shanraq.com/internal/web"
)

// Router exposes the public-facing endpoints.
func Router(cfg config.Config, logger zerolog.Logger, renderer *web.Renderer) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("path", r.URL.Path).Msg("public_page")
		w.Header().Set("X-App-Name", cfg.App.Name)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if renderer != nil {
			data := web.HomePageData{
				BrandName:   cfg.App.Name,
				PageTitle:   "Global Real Estate Platform · ",
				Description: "Discover, list, and manage properties and logistics partners across the world with Shanraq.",
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
