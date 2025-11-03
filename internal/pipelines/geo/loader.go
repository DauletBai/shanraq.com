package geo

import (
	"context"
	"database/sql"
	"time"
)

// Loader represents a pipeline stage capable of ingesting geo data.
type Loader struct {
	db *sql.DB
}

// NewLoader builds a new geo pipeline loader.
func NewLoader(db *sql.DB) *Loader {
	return &Loader{db: db}
}

// Run performs a single ingestion cycle. In demo form it is a no-op with observability hooks.
func (l *Loader) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Millisecond):
		// In a future iteration we will call external geo APIs and upsert into the database.
		return nil
	}
}
