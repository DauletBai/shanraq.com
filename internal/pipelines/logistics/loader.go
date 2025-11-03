package logistics

import (
	"context"
	"database/sql"
	"time"
)

// Loader orchestrates ingestion of logistics partner updates.
type Loader struct {
	db *sql.DB
}

func NewLoader(db *sql.DB) *Loader {
	return &Loader{db: db}
}

// Run executes the ingestion pipeline. For now it simulates processing latency.
func (l *Loader) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(15 * time.Millisecond):
		return nil
	}
}
