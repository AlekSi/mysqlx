package mysqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net"
)

type OpenParams struct {
	Dialer *net.Dialer
	Trace  TraceFunc
}

var defaultOpenParams = &OpenParams{
	Dialer: new(net.Dialer),
	Trace:  noTrace,
}

// Driver implements database/sql/driver.Driver interface.
type Driver struct{}

// Open returns a new connection to the database. See README for dataSource format.
// The returned connection may be used only by one goroutine at a time.
func (Driver) Open(dataSource string) (driver.Conn, error) {
	return open(context.Background(), dataSource, defaultOpenParams)
}

// OpenWithParams is a variant of Open with specified context and other parameters.
// Context is used only for connection establishing: dialing (with params.Dialer.DialContext), negotiating
// and authenticating. Canceling the context after the connection is established does nothing.
func (Driver) OpenWithParams(ctx context.Context, dataSource string, params *OpenParams) (driver.Conn, error) {
	return open(ctx, dataSource, params)
}

func init() {
	sql.Register("mysqlx", Driver{})
}

// check interfaces
var (
	_ driver.Driver = Driver{}
)
