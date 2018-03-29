// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"context"
	"database/sql/driver"
)

// tx represents transaction.
//
// It implements driver.Tx.
type tx struct {
	c *conn
}

// Commit commits the transaction.
func (t *tx) Commit() error {
	_, err := t.c.ExecContext(context.Background(), "COMMIT", nil)
	return err
}

// Rollback aborts the transaction.
func (t *tx) Rollback() error {
	_, err := t.c.ExecContext(context.Background(), "ROLLBACK", nil)
	return err
}

// check interfaces
var (
	_ driver.Tx = (*tx)(nil)
)
