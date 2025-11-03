package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"shanraq.com/internal/auth"
	"shanraq.com/internal/auth/session"
	"shanraq.com/internal/config"
	"shanraq.com/internal/database"
	"shanraq.com/internal/httpserver"
	"shanraq.com/internal/logging"
	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
	workspaceservice "shanraq.com/internal/services/workspace"
	"shanraq.com/internal/web"
)

// App encapsulates the primary application lifecycle.
type App struct {
	cfg          config.Config
	logger       zerolog.Logger
	server       *httpserver.Server
	db           *sql.DB
	transportSvc transportservice.Service
	agencySvc    agencyservice.Service
	listingSvc   listingservice.Service
	authRegistry *auth.ProviderRegistry
	sessions     *session.Manager
	workspaces   workspaceservice.Service
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

	var transportSvc transportservice.Service = transportservice.NewInMemoryService()
	var agencySvc agencyservice.Service = agencyservice.NewInMemoryService()
	var listingSvc listingservice.Service = listingservice.NewInMemoryService()
	var workspaceSvc workspaceservice.Service = workspaceservice.NewInMemoryService()

	var db *sql.DB
	if cfg.Database.URL != "" {
		if conn, err := database.Connect(context.Background(), cfg.Database); err != nil {
			logger.Warn().Err(err).Msg("database connection failed; using in-memory services")
		} else {
			db = conn
			if svc, err := transportservice.NewSQLService(conn); err != nil {
				logger.Warn().Err(err).Msg("init transport sql service")
			} else {
				transportSvc = svc
			}
			if svc, err := agencyservice.NewSQLService(conn); err != nil {
				logger.Warn().Err(err).Msg("init agency sql service")
			} else {
				agencySvc = svc
			}
			if svc, err := listingservice.NewSQLService(conn); err != nil {
				logger.Warn().Err(err).Msg("init listing sql service")
			} else {
				listingSvc = svc
			}
		}
	}
	authRegistry := auth.NewRegistry(cfg.Auth.SupportedProviders...)
	for _, name := range cfg.Auth.SupportedProviders {
		authRegistry.Register(name, auth.NewDemoOAuthProvider(name))
	}
	if cfg.Auth.Provider != "" {
		authRegistry.Register(cfg.Auth.Provider, auth.NewDemoOAuthProvider(cfg.Auth.Provider))
	}
	sessionManager := session.NewManager(12*time.Hour, "")

	router := httpserver.NewRouter(httpserver.Deps{
		Logger:           logger,
		Config:           cfg,
		Renderer:         renderer,
		TransportService: transportSvc,
		AgencyService:    agencySvc,
		ListingService:   listingSvc,
		AuthRegistry:     authRegistry,
		SessionManager:   sessionManager,
		WorkspaceService: workspaceSvc,
	})

	server := httpserver.New(cfg.HTTP, router, logger)

	return &App{
		cfg:          cfg,
		logger:       logger,
		server:       server,
		db:           db,
		transportSvc: transportSvc,
		agencySvc:    agencySvc,
		listingSvc:   listingSvc,
		authRegistry: authRegistry,
		sessions:     sessionManager,
		workspaces:   workspaceSvc,
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
		if a.db != nil {
			if err := a.db.Close(); err != nil {
				a.logger.Warn().Err(err).Msg("close database")
			}
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
