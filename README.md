# Shanraq.com Platform

Shanraq is an international real-estate platform that connects property owners, buyers, realtors, agencies, corporate landlords, and country-specific logistics providers under a single, AI-augmented experience. The project is built with Go, Chi, PostgreSQL, and a modern frontend stack focused on global reach, localization, and extensibility.

## Highlights
- Go (1.22+) backend with Chi, structured logging, and modular internal packages.
- Blog-inspired public landing page served from HTML templates with Bootstrap 5 styling.
- REST foundation for upcoming CRUD modules, including transportation company onboarding, agency/realtor directories, and property listings.
- Configurable via environment variables (`APP_ENV`, `DATABASE_URL`, etc.) using `envconfig`.
- Static assets and dashboard prototype mounted directly from the `web` directory.
- Multi-provider auth facade prepared for Google, Meta, Apple, LinkedIn, e-mail magic links, and a primary provider (`AUTH_PROVIDER`).
- Landing page renders seeded demo data (agencies, realtors, moving partners, and 10 international listings) for quick UX validation.
- Server-side templating powered by Go `html/template` with a shared layout (`web/layout.html`), modular partials, and per-page views under `web/pages/` for easy expansion.

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

### Makefile Shortcuts

```bash
# build binaries under bin/
make build

# start the API (same as go run ./cmd/app)
make run

# tidy modules, vet, and test
make tidy lint test

# build migration helper and apply migrations
make migrate-up DATABASE_URL=postgres://...

# roll back one migration
make migrate-down DATABASE_URL=postgres://...
```

## Demo Data & APIs

- `GET /api/v1/listings` — seeded international properties (10 items) with `/featured` variant for homepage cards.
- `GET /api/v1/agencies` — global agencies with `/realtors` and `/realtors/featured`.
- `GET /api/v1/transport-companies` — moving/logistics partners with regional coverage metadata.
- `GET /auth/providers` — lists configured authentication providers (Google, Meta, Apple, LinkedIn, Email, plus primary provider).
- Landing page consumes the same demo data to showcase cards for listings, agencies, realtors, and logistics firms.

Environment variables allow configuration:

| Variable | Description | Default |
| --- | --- | --- |
| `APP_ENV` | Runtime environment (`development`, `production`, …) | `development` |
| `HTTP_PORT` | HTTP port binding | `8080` |
| `DATABASE_URL` | PostgreSQL connection string (PostGIS-ready) | `postgres://postgres:postgres@localhost:5432/shanraq?sslmode=disable` |
| `FEATURES_ENABLE_TRANSPORT_COMPANIES` | Toggle transportation company workflows | `true` |
| `AUTH_SUPPORTED_PROVIDERS` | Comma-separated list of login providers exposed at `/auth/providers` | `google,meta,linkedin,apple,email` |

## CI & Branch Protection

- GitHub Actions workflow **Go CI** (see `.github/workflows/ci.yml`) runs `gofmt`, `go vet`, `go build`, and `go test` on every push and pull request targeting `main`.
- Before inviting collaborators, protect the `main` branch in GitHub settings:
  1. Settings → Branches → Add rule for `main`.
  2. Require pull requests before merging, with at least one approval.
  3. Check “Require status checks to pass” and select **Go CI**.
  4. Enable “Require branches to be up to date” and “Include administrators” as fits your governance.
- Optionally enforce signed commits and restrict who can push directly to `main`.


## Next Steps

- **Persistence layer**: introduce PostgreSQL migrations (Atlas or golang-migrate), domain repositories, and seeders for countries/regions/cities.
- **Authentication**: integrate OpenID Connect/social login and role-based access (buyers, realtors, agencies, logistics partners, admins).
- **Transport companies**: extend in-memory CRUD to PostgreSQL, add compliance workflows, and availability per country (especially markets with strong logistics demand).
- **Agencies & Listings**: back services with database repositories, add search filters, media galleries, and multilingual content.
- **Search & AI**: add PostGIS/Elasticsearch for geo + full-text, and LLM-powered content (translated listings, recommendations, moderation).
- **Dashboard**: evolve the `/dashboard` UI into an interactive console for inventory, analytics, and operations.
- **Ops**: ship Docker Compose, CI/CD, telemetry (OpenTelemetry + Prometheus), and security hardening (TLS, secrets management).

Contributions follow the Go project layout already established (`cmd/`, `internal/`, `pkg/`). Each module should include targeted unit/integration testing and documentation. Continuous collaboration on UX copy, localization, and data ingestion pipelines will ensure Shanraq stands out as the best global real estate experience.
