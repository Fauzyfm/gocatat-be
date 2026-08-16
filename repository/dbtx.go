// repository/dbtx.go — file baru, shared interface
package repository

import (
    "context"
    "database/sql"
)

// DBTX memungkinkan method menerima *sql.DB atau *sql.Tx
// sehingga bisa dipakai di dalam maupun di luar DB transaction
type DBTX interface {
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}