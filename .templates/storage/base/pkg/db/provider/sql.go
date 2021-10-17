package provider

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// SQL implements default sql provider.
type SQL struct {
	mutex sync.RWMutex
	db    *sql.DB
	tx    *sql.Tx
	ctx   context.Context
}

// New returns new sql provider.
func New(db *sql.DB) *SQL {
	return &SQL{db: db}
}

// TransactProvider returns new sql provider with transaction.
func (s *SQL) TransactProvider() (*SQL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	tx, err := s.db.Begin()

	if err != nil {
		return s, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &SQL{db: s.db, tx: tx, ctx: s.ctx}, nil
}

// Commit changes in depth of transaction.
func (s *SQL) Commit() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.tx != nil {
		defer func() { s.tx = nil }()

		if err := s.tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit changes %w", err)
		}

		return nil
	}

	return ErrNotDefinedTransaction
}

// Rollback changes in depth of transaction.
func (s *SQL) Rollback() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.tx != nil {
		defer func() { s.tx = nil }()

		if err := s.tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback changes %w", err)
		}

		return nil
	}

	return ErrNotDefinedTransaction
}

// Context returns sql provider with context.
func (s *SQL) Context(ctx context.Context) *SQL {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return &SQL{db: s.db, tx: s.tx, ctx: ctx}
}

// Query does sql request and returns rows.
func (s *SQL) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	switch {
	case s.tx != nil && s.ctx != nil:
		rows, err = s.tx.QueryContext(s.ctx, query, args...)
	case s.tx != nil && s.ctx == nil:
		rows, err = s.tx.Query(query, args...)
	case s.tx == nil && s.ctx != nil:
		rows, err = s.db.QueryContext(s.ctx, query, args...)
	default:
		rows, err = s.db.Query(query, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query %w", err)
	}

	return
}

// QueryRow does sql request and returns row.
func (s *SQL) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.tx == nil {
		if s.ctx == nil {
			return s.db.QueryRow(query, args...)
		}

		return s.db.QueryRowContext(s.ctx, query, args...)
	}

	if s.ctx == nil {
		return s.tx.QueryRow(query, args...)
	}

	return s.tx.QueryRowContext(s.ctx, query, args...)
}

// Prepare does sql request and return sql statement.
func (s *SQL) Prepare(query string) (stmt *sql.Stmt, err error) {
	switch {
	case s.tx != nil && s.ctx != nil:
		stmt, err = s.tx.PrepareContext(s.ctx, query)
	case s.tx != nil && s.ctx == nil:
		stmt, err = s.tx.Prepare(query)
	case s.tx == nil && s.ctx != nil:
		stmt, err = s.db.PrepareContext(s.ctx, query)
	default:
		stmt, err = s.db.Prepare(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to prepare query %w", err)
	}

	return
}
