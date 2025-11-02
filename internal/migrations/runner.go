package migrations

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Runner executes filesystem-based migrations against PostgreSQL.
type Runner struct {
	databaseURL string
	dir         string
}

// NewRunner builds a Runner instance.
func NewRunner(databaseURL, dir string) *Runner {
	return &Runner{databaseURL: databaseURL, dir: dir}
}

// Up applies all available migrations.
func (r *Runner) Up() error {
	m, err := r.migrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Down reverts the most recent `steps` migrations.
func (r *Runner) Down(steps int) error {
	m, err := r.migrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if steps <= 0 {
		steps = 1
	}

	return m.Steps(-steps)
}

// Steps executes an arbitrary number of steps (positive or negative).
func (r *Runner) Steps(steps int) error {
	m, err := r.migrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if steps == 0 {
		return nil
	}

	return m.Steps(steps)
}

func (r *Runner) migrator() (*migrate.Migrate, error) {
	absDir, err := filepath.Abs(r.dir)
	if err != nil {
		return nil, fmt.Errorf("resolve migrations directory: %w", err)
	}
	uri := fmt.Sprintf("file://%s", absDir)
	return migrate.New(uri, r.databaseURL)
}
