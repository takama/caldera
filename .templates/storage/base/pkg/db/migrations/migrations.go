package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

// Migrator design migration interface
type Migrator interface {
	Setup(db *sql.DB)
	Migrate() error
	MigrateUp(version int64) error
	MigrateDown(version int64) error
}

// Sequence contains migration functionality
type Sequence struct {
	db  *sql.DB
	cfg *Config
}

// New creates migration Sequence
func New(cfg *Config) *Sequence {
	return &Sequence{
		cfg: cfg,
	}
}

// Setup configure migration
func (s *Sequence) Setup(db *sql.DB) {
	s.db = db
}

// Migrate process migration up to last version
func (s Sequence) Migrate() error {
	if s.cfg.Active && s.db != nil {
		return goose.Up(s.db, s.cfg.Dir)
	}

	return nil
}

// MigrateUp process migration up to specified version
// Use version 0 if plan migrate up to last version
func (s Sequence) MigrateUp(version int64) error {
	if s.cfg.Active && s.db != nil {
		if version == 0 {
			return goose.Up(s.db, s.cfg.Dir)
		}
		return goose.UpTo(s.db, s.cfg.Dir, version)
	}

	return nil
}

// MigrateDown process migration down to specified version
// Use version 0 if plan migrate down to initial state
func (s Sequence) MigrateDown(version int64) error {
	if s.cfg.Active && s.db != nil {
		if version == 0 {
			return goose.Down(s.db, s.cfg.Dir)
		}
		return goose.DownTo(s.db, s.cfg.Dir, version)
	}

	return nil
}
