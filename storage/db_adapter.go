package storage

import (
	"context"
	"database/sql"
)

// DBAdapter wraps sql.DB to implement DBInterface
type DBAdapter struct {
	*sql.DB
}

// NewDBAdapter creates a new DBAdapter
func NewDBAdapter(db *sql.DB) *DBAdapter {
	return &DBAdapter{DB: db}
}

// QueryRowContext implements DBInterface.QueryRowContext
func (a *DBAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) RowInterface {
	row := a.DB.QueryRowContext(ctx, query, args...)
	return &RowAdapter{Row: row}
}

// RowAdapter wraps sql.Row to implement RowInterface
type RowAdapter struct {
	*sql.Row
}

// Scan implements RowInterface.Scan
func (r *RowAdapter) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}
