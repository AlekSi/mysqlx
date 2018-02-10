// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net/url"
)

// driverType implements database/sql/driver.Driver interface.
type driverType struct{}

var Driver driverType

// Open returns a new connection to the database. See README for dataSource format.
// The returned connection may be used only by one goroutine at a time.
func (driverType) Open(dataSource string) (driver.Conn, error) {
	u, err := url.Parse(dataSource)
	if err != nil {
		return nil, err
	}
	ds, err := ParseDataSource(u)
	if err != nil {
		return nil, err
	}
	return open(context.Background(), ds)
}

/*
// OpenCtx is a variant of Open with specified context and DataSource struct.
// Context is used only for connection establishing: dialing, negotiating
// and authenticating. Canceling the context after the connection is established does nothing.
func (driverType) OpenCtx(ctx context.Context, dataSource *DataSource) (driver.Conn, error) {
	return open(ctx, dataSource)
}
*/

// TODO
// func (Driver) OpenConnector(name string) (driver.Connector, error) {
// 	return nil, nil
// }

func init() {
	sql.Register("mysqlx", Driver)
}

// check interfaces
var (
	_ driver.Driver = Driver
	// TODO _ driver.DriverContext = Driver
)
