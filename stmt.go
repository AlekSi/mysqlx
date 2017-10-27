package mysqlx

import (
	"context"
	"database/sql/driver"
)

// stmt is a prepared statement. It is bound to a Conn and not used by multiple goroutines concurrently.
type stmt struct {
	c     *conn
	query string
}

// Close closes the statement.
func (s *stmt) Close() error {
	return nil
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (s *stmt) NumInput() int {
	return -1
}

// Exec executes a query that doesn't return rows, such as an INSERT or UPDATE.
func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.c.Exec(s.query, args)
}

// ExecContext executes a query that doesn't return rows, such as an INSERT or UPDATE.
// It honors the context timeout and return when it is canceled.
func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.c.ExecContext(ctx, s.query, args)
}

// Query executes a query that may return rows, such as a SELECT.
func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.c.Query(s.query, args)
}

// Query executes a query that may return rows, such as a SELECT.
// It honors the context timeout and return when it is canceled.
func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return s.c.QueryContext(ctx, s.query, args)
}

// check interfaces
var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)

	// TODO
	// _ driver.ColumnConverter = (*stmt)(nil)
	// _ driver.NamedValueChecker = (*stmt)(nil)
)
