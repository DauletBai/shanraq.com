# Database Migrations

The project uses [golang-migrate](https://github.com/golang-migrate/migrate) compatible SQL migrations. Each migration is versioned using the `<version>_<name>.up.sql`/`.down.sql` naming convention.

## Running Locally

1. Build the migration helper:

   ```bash
   go build -o bin/migrate ./cmd/cli/migrate
   ```

2. Execute migrations against your database:

   ```bash
   bin/migrate \
     -database "$DATABASE_URL" \
     -dir migrations \
     up
   ```

   The `DATABASE_URL` should follow the standard PostgreSQL connection format, e.g.:
   `postgres://postgres:postgres@localhost:5432/shanraq?sslmode=disable`

3. To roll back the most recent migration:

   ```bash
   bin/migrate \
     -database "$DATABASE_URL" \
     -dir migrations \
     down 1
   ```

## Tips

- Ensure the `pgcrypto` extension is available; the first migration enables it to support UUID primary keys.
- For local development you can use Docker Compose to provision PostgreSQL + PostGIS (planned in `/deploy`).
- Keep migrations idempotent and include corresponding `down` scripts when possible.
