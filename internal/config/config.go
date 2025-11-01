package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Config aggregates runtime configuration for the application.
	Config struct {
		App        App        `envconfig:"APP"`
		HTTP       HTTP       `envconfig:"HTTP"`
		Database   Database   `envconfig:"DATABASE"`
		Telemetry  Telemetry  `envconfig:"TELEMETRY"`
		Auth       Auth       `envconfig:"AUTH"`
		Geo        Geo        `envconfig:"GEO"`
		AI         AI         `envconfig:"AI"`
		Features   Features   `envconfig:"FEATURES"`
		Seed       Seed       `envconfig:"SEED"`
		Scheduling Scheduling `envconfig:"SCHEDULING"`
	}

	App struct {
		Name    string `envconfig:"NAME" default:"shanraq"`
		Env     string `envconfig:"ENV" default:"development"`
		Version string `envconfig:"VERSION" default:"0.1.0"`
	}

	HTTP struct {
		Host           string        `envconfig:"HOST" default:"0.0.0.0"`
		Port           int           `envconfig:"PORT" default:"8080"`
		ReadTimeout    time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout   time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
		IdleTimeout    time.Duration `envconfig:"IDLE_TIMEOUT" default:"120s"`
		PublicBaseURL  string        `envconfig:"PUBLIC_BASE_URL" default:"http://localhost:8080"`
		DashboardURL   string        `envconfig:"DASHBOARD_URL" default:"http://localhost:8080/dashboard"`
		AllowedOrigins []string      `envconfig:"ALLOWED_ORIGINS" default:"*"`
	}

	Database struct {
		URL               string        `envconfig:"URL" default:"postgres://postgres:postgres@localhost:5432/shanraq?sslmode=disable"`
		MaxOpenConn       int           `envconfig:"MAX_OPEN_CONN" default:"25"`
		MaxIdleConn       int           `envconfig:"MAX_IDLE_CONN" default:"25"`
		ConnMaxIdleTime   time.Duration `envconfig:"CONN_MAX_IDLE_TIME" default:"15m"`
		ConnMaxLifetime   time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"60m"`
		MigrationDir      string        `envconfig:"MIGRATION_DIR" default:"migrations"`
		GeoSeedDataSource string        `envconfig:"GEO_SEED_DATA_SOURCE" default:"data/geo"`
	}

	Telemetry struct {
		Enabled      bool    `envconfig:"ENABLED" default:"false"`
		ServiceName  string  `envconfig:"SERVICE_NAME" default:"shanraq-api"`
		ExporterURL  string  `envconfig:"EXPORTER_URL" default:"http://localhost:4317"`
		MetricsBind  string  `envconfig:"METRICS_BIND" default:":9464"`
		TracingRatio float64 `envconfig:"TRACING_RATIO" default:"0.1"`
	}

	Auth struct {
		Provider       string   `envconfig:"PROVIDER" default:"auth0"`
		ClientID       string   `envconfig:"CLIENT_ID"`
		ClientSecret   string   `envconfig:"CLIENT_SECRET"`
		CallbackURL    string   `envconfig:"CALLBACK_URL" default:"http://localhost:8080/auth/callback"`
		AllowedDomains []string `envconfig:"ALLOWED_DOMAINS"`
		JWTSigningKey  string   `envconfig:"JWT_SIGNING_KEY" default:"development-secret"`
	}

	Geo struct {
		DataProvider string        `envconfig:"DATA_PROVIDER" default:"geonames"`
		CacheTTL     time.Duration `envconfig:"CACHE_TTL" default:"6h"`
	}

	AI struct {
		Provider         string  `envconfig:"PROVIDER" default:"openai"`
		Endpoint         string  `envconfig:"ENDPOINT"`
		APIKey           string  `envconfig:"API_KEY"`
		EnableListings   bool    `envconfig:"ENABLE_LISTINGS" default:"true"`
		EnableModeration bool    `envconfig:"ENABLE_MODERATION" default:"true"`
		BudgetUSD        float64 `envconfig:"BUDGET_USD" default:"50"`
	}

	Features struct {
		EnableDashboard          bool `envconfig:"ENABLE_DASHBOARD" default:"true"`
		EnableTransportCompanies bool `envconfig:"ENABLE_TRANSPORT_COMPANIES" default:"true"`
		EnableAgencies           bool `envconfig:"ENABLE_AGENCIES" default:"true"`
		EnableAIRecommendations  bool `envconfig:"ENABLE_AI_RECOMMENDATIONS" default:"true"`
	}

	Seed struct {
		EnableAutoSeed bool          `envconfig:"ENABLE_AUTO_SEED" default:"true"`
		ChunkSize      int           `envconfig:"CHUNK_SIZE" default:"500"`
		Timeout        time.Duration `envconfig:"TIMEOUT" default:"10m"`
		RegionsFilter  []string      `envconfig:"REGIONS_FILTER"`
	}

	Scheduling struct {
		EnableJobs bool          `envconfig:"ENABLE_JOBS" default:"true"`
		Timezone   string        `envconfig:"TIMEZONE" default:"UTC"`
		Interval   time.Duration `envconfig:"INTERVAL" default:"1m"`
	}
)

// Load reads configuration from environment variables.
func Load() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
