package mysqlx

import (
	"database/sql/driver"
)

type tx struct {
	c *conn
}

func (t *tx) Commit() error {
	_, err := t.c.Exec("COMMIT", nil)
	return err
}

func (t *tx) Rollback() error {
	_, err := t.c.Exec("ROLLBACK", nil)
	return err
}

// check interfaces
var (
	_ driver.Tx = (*tx)(nil)
)
