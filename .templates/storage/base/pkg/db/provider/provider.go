package provider

import (
	"context"
	"database/sql"
	"sync"
)

// Provider implements default sql provider
type Provider struct {
	mutex sync.RWMutex
	db    *sql.DB
	tx    *sql.Tx
	ctx   context.Context
}

// New returns new sql provider
func New(db *sql.DB) *Provider {
	return &Provider{db: db}
}

// TransactProvider returns new provider with transaction
func (p *Provider) TransactProvider() (*Provider, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	tx, err := p.db.Begin()
	if err != nil {
		return p, err
	}
	return &Provider{db: p.db, tx: tx, ctx: p.ctx}, nil
}

// Commit changes in depth of transaction
func (p *Provider) Commit() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.tx != nil {
		defer func() { p.tx = nil }()
		return p.tx.Commit()
	}
	return ErrNotDefinedTransaction
}

// Rollback changes in depth of transaction
func (p *Provider) Rollback() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.tx != nil {
		defer func() { p.tx = nil }()
		return p.tx.Rollback()
	}
	return ErrNotDefinedTransaction
}

// Context returns provider with context
func (p *Provider) Context(ctx context.Context) *Provider {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return &Provider{db: p.db, tx: p.tx, ctx: ctx}
}

// Query does sql request and returns rows
func (p *Provider) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if p.tx == nil {
		if p.ctx == nil {
			rows, err = p.db.Query(query, args...)
		} else {
			rows, err = p.db.QueryContext(p.ctx, query, args...)
		}
	} else {
		if p.ctx == nil {
			rows, err = p.tx.Query(query, args...)
		} else {
			rows, err = p.tx.QueryContext(p.ctx, query, args...)
		}
	}

	return rows, err
}

// QueryRow does sql request and returns row
func (p *Provider) QueryRow(query string, args ...interface{}) *sql.Row {
	var row *sql.Row
	if p.tx == nil {
		if p.ctx == nil {
			row = p.db.QueryRow(query, args...)
		} else {
			row = p.db.QueryRowContext(p.ctx, query, args...)
		}
	} else {
		if p.ctx == nil {
			row = p.tx.QueryRow(query, args...)
		} else {
			row = p.tx.QueryRowContext(p.ctx, query, args...)
		}
	}

	return row
}

// Prepare does sql request and return sql statement
func (p *Provider) Prepare(query string) (*sql.Stmt, error) {
	if p.tx == nil {
		if p.ctx == nil {
			return p.db.Prepare(query)
		}
		return p.db.PrepareContext(p.ctx, query)
	}
	if p.ctx == nil {
		return p.tx.Prepare(query)
	}
	return p.tx.PrepareContext(p.ctx, query)
}
