package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

const errApplyMigration = "failed to apply mifration %w"

// Migrator design migration interface.
type Migrator interface {
	Setup(db *sql.DB) error
	Migrate() error
	MigrateUp(version int64) error
	MigrateDown(version int64) error
}

// Sequence contains migration functionality.
type Sequence struct {
	db  *sql.DB
	cfg *Config
}

// New creates migration Sequence.
func New(cfg *Config) *Sequence {
	return &Sequence{
		cfg: cfg,
	}
}

// Setup configure migration.
func (s *Sequence) Setup(db *sql.DB) error {
	s.db = db

	if err := goose.SetDialect(s.cfg.Dialect); err != nil {
		return fmt.Errorf("failed to set migration dialect %w", err)
	}

	return nil
}

// Migrate process migration up to last version.
func (s Sequence) Migrate() error {
	if s.cfg.Active && s.db != nil {
		if err := goose.Up(s.db, s.cfg.Dir); err != nil {
			return fmt.Errorf(errApplyMigration, err)
		}
	}

	return nil
}

// MigrateUp process migration up to specified version.
// Use version 0 if plan migrate up to last version.
func (s Sequence) MigrateUp(version int64) error {
	if s.cfg.Active && s.db != nil {
		if version == 0 {
			if err := goose.Up(s.db, s.cfg.Dir); err != nil {
				return fmt.Errorf(errApplyMigration, err)
			}

			return nil
		}

		if err := goose.UpTo(s.db, s.cfg.Dir, version); err != nil {
			return fmt.Errorf(errApplyMigration, err)
		}
	}

	return nil
}

// MigrateDown process migration down to specified version.
// Use version 0 if plan migrate down to initial state.
func (s Sequence) MigrateDown(version int64) error {
	if s.cfg.Active && s.db != nil {
		if version == 0 {
			if err := goose.Down(s.db, s.cfg.Dir); err != nil {
				return fmt.Errorf(errApplyMigration, err)
			}

			return nil
		}

		if err := goose.DownTo(s.db, s.cfg.Dir, version); err != nil {
			return fmt.Errorf(errApplyMigration, err)
		}
	}

	return nil
}
