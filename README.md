# Shanraq.com Platform

Shanraq is an international real-estate platform that connects property owners, buyers, realtors, agencies, corporate landlords, and country-specific logistics providers under a single, AI-augmented experience. The project is built with Go, Chi, PostgreSQL, and a modern frontend stack focused on global reach, localization, and extensibility.

## Highlights
- Go (1.22+) backend with Chi, structured logging, and modular internal packages.
- Blog-inspired public landing page served from HTML templates with Bootstrap 5 styling.
- REST foundation for upcoming CRUD modules, including transportation company onboarding.
- Configurable via environment variables (`APP_ENV`, `DATABASE_URL`, etc.) using `envconfig`.
- Static assets and dashboard prototype mounted directly from the `web` directory.

## Project Layout

```
├── cmd/
│   ├── app/        # HTTP server entry point
│   └── cli/        # Future CLI utilities (migrations, seeding, AI tooling)
├── internal/
│   ├── app/        # Application lifecycle wiring
│   ├── config/     # Environment-driven configuration structs
│   ├── httpserver/ # Router, middleware, API/public handlers
│   ├── logging/    # Zerolog helpers
│   └── web/        # Template renderer, helper functions
├── pkg/            # Shared libraries (reserved)
├── web/            # Templates, dashboard prototype, static assets
└── README.md
```

## Getting Started

1. **Install dependencies** (Go 1.22+).
2. **Run the server**:

   ```bash
   go run ./cmd/app
   ```

   The public site is available at `http://localhost:8080` and the dashboard prototype at `http://localhost:8080/dashboard/`.

3. **Health check**: `GET http://localhost:8080/healthz` returns `200 OK`.

Environment variables allow configuration:

| Variable | Description | Default |
| --- | --- | --- |
| `APP_ENV` | Runtime environment (`development`, `production`, …) | `development` |
| `HTTP_PORT` | HTTP port binding | `8080` |
| `DATABASE_URL` | PostgreSQL connection string (PostGIS-ready) | `postgres://postgres:postgres@localhost:5432/shanraq?sslmode=disable` |
| `FEATURES_ENABLE_TRANSPORT_COMPANIES` | Toggle transportation company workflows | `true` |

## Next Steps

- **Persistence layer**: introduce PostgreSQL migrations (Atlas or golang-migrate), domain repositories, and seeders for countries/regions/cities.
- **Authentication**: integrate OpenID Connect/social login and role-based access (buyers, realtors, agencies, logistics partners, admins).
- **Transport companies**: implement full CRUD, compliance workflows, and availability per country (especially markets with strong logistics demand).
- **Search & AI**: add PostGIS/Elasticsearch for geo + full-text, and LLM-powered content (translated listings, recommendations, moderation).
- **Dashboard**: evolve the `/dashboard` UI into an interactive console for inventory, analytics, and operations.
- **Ops**: ship Docker Compose, CI/CD, telemetry (OpenTelemetry + Prometheus), and security hardening (TLS, secrets management).

Contributions follow the Go project layout already established (`cmd/`, `internal/`, `pkg/`). Each module should include targeted unit/integration testing and documentation. Continuous collaboration on UX copy, localization, and data ingestion pipelines will ensure Shanraq stands out as the best global real estate experience.
