package mysqlx

import (
	"database/sql"
	"database/sql/driver"
)

type mysqlxDriver struct{}

func (mysqlxDriver) Open(dataSource string) (driver.Conn, error) {
	return open(dataSource)
}

func init() {
	sql.Register("mysqlx", mysqlxDriver{})
}

// check interfaces
var (
	_ driver.Driver = mysqlxDriver{}
)
