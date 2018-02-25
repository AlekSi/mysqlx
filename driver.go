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
)

// driverType implements database/sql/driver.Driver interface.
type driverType struct{}

var Driver driverType

// Open returns a new connection to the database. See README for dataSource format.
// The returned connection may be used only by one goroutine at a time.
func (d driverType) Open(dataSource string) (driver.Conn, error) {
	ds, err := d.OpenConnector(dataSource)
	if err != nil {
		return nil, err
	}
	return open(context.Background(), ds.(*DataSource))
}

func (d driverType) OpenConnector(dataSource string) (driver.Connector, error) {
	ds, err := ParseDataSource(dataSource)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func init() {
	sql.Register("mysqlx", Driver)
}

// check interfaces
var (
	_ driver.Driver        = Driver
	_ driver.DriverContext = Driver
)
