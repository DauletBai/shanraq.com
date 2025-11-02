package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"shanraq.com/internal/auth"
	"shanraq.com/internal/config"
	"shanraq.com/internal/httpserver"
	"shanraq.com/internal/logging"
	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
	"shanraq.com/internal/web"
)

// App encapsulates the primary application lifecycle.
type App struct {
	cfg          config.Config
	logger       zerolog.Logger
	server       *httpserver.Server
	transportSvc transportservice.Service
	agencySvc    agencyservice.Service
	listingSvc   listingservice.Service
	authRegistry *auth.ProviderRegistry
}

// New wires the core application dependencies.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	normalizeConfig(&cfg)

	logger := logging.New(cfg.App.Env)

	renderer, err := web.NewRenderer()
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	transportSvc := transportservice.NewInMemoryService()
	agencySvc := agencyservice.NewInMemoryService()
	listingSvc := listingservice.NewInMemoryService()
	authRegistry := auth.NewRegistry(cfg.Auth.SupportedProviders...)
	if cfg.Auth.Provider != "" {
		authRegistry.Register(cfg.Auth.Provider, auth.NewNoopProvider())
	}

	router := httpserver.NewRouter(httpserver.Deps{
		Logger:           logger,
		Config:           cfg,
		Renderer:         renderer,
		TransportService: transportSvc,
		AgencyService:    agencySvc,
		ListingService:   listingSvc,
		AuthRegistry:     authRegistry,
	})

	server := httpserver.New(cfg.HTTP, router, logger)

	return &App{
		cfg:          cfg,
		logger:       logger,
		server:       server,
		transportSvc: transportSvc,
		agencySvc:    agencySvc,
		listingSvc:   listingSvc,
		authRegistry: authRegistry,
	}, nil
}

// Run starts the HTTP server and blocks until context cancellation or fatal error.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := a.server.Start(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info().Msg("shutdown signal received")
		if err := a.server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}
		return nil

	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("http server: %w", err)
	}
}

func normalizeConfig(cfg *config.Config) {
	if len(cfg.HTTP.AllowedOrigins) == 1 && cfg.HTTP.AllowedOrigins[0] == "" {
		cfg.HTTP.AllowedOrigins = nil
	}
	if len(cfg.HTTP.AllowedOrigins) == 0 {
		cfg.HTTP.AllowedOrigins = []string{"*"}
	}
}
