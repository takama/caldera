package provider

import (
	"context"
	"database/sql"
	"sync"
)

// SQLProvider implements default sql provider
type SQLProvider struct {
	mutex sync.RWMutex
	db    *sql.DB
	tx    *sql.Tx
	ctx   context.Context
}

// New returns new sql provider
func New(db *sql.DB) *SQLProvider {
	return &SQLProvider{db: db}
}

// TransactProvider returns new provider with transaction
func (sp *SQLProvider) TransactProvider() (*SQLProvider, error) {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()
	tx, err := sp.db.Begin()
	if err != nil {
		return sp, err
	}
	return &SQLProvider{db: sp.db, tx: tx, ctx: sp.ctx}, nil
}

// Commit changes in depth of transaction
func (sp *SQLProvider) Commit() error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	if sp.tx != nil {
		defer func() { sp.tx = nil }()
		return sp.tx.Commit()
	}
	return ErrNotDefinedTransaction
}

// Rollback changes in depth of transaction
func (sp *SQLProvider) Rollback() error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	if sp.tx != nil {
		defer func() { sp.tx = nil }()
		return sp.tx.Rollback()
	}
	return ErrNotDefinedTransaction
}

// Context returns provider with context
func (sp *SQLProvider) Context(ctx context.Context) *SQLProvider {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()
	return &SQLProvider{db: sp.db, tx: sp.tx, ctx: ctx}
}

// Query does sql request and returns rows
func (sp *SQLProvider) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if sp.tx == nil {
		if sp.ctx == nil {
			rows, err = sp.db.Query(query, args...)
		} else {
			rows, err = sp.db.QueryContext(sp.ctx, query, args...)
		}
	} else {
		if sp.ctx == nil {
			rows, err = sp.tx.Query(query, args...)
		} else {
			rows, err = sp.tx.QueryContext(sp.ctx, query, args...)
		}
	}

	return rows, err
}

// QueryRow does sql request and returns row
func (sp *SQLProvider) QueryRow(query string, args ...interface{}) *sql.Row {
	var row *sql.Row
	if sp.tx == nil {
		if sp.ctx == nil {
			row = sp.db.QueryRow(query, args...)
		} else {
			row = sp.db.QueryRowContext(sp.ctx, query, args...)
		}
	} else {
		if sp.ctx == nil {
			row = sp.tx.QueryRow(query, args...)
		} else {
			row = sp.tx.QueryRowContext(sp.ctx, query, args...)
		}
	}

	return row
}

// Prepare does sql request and return sql statement
func (sp *SQLProvider) Prepare(query string) (*sql.Stmt, error) {
	if sp.tx == nil {
		if sp.ctx == nil {
			return sp.db.Prepare(query)
		}
		return sp.db.PrepareContext(sp.ctx, query)
	}
	if sp.ctx == nil {
		return sp.tx.Prepare(query)
	}
	return sp.tx.PrepareContext(sp.ctx, query)
}
