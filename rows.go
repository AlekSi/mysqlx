package mysqlx

import (
	"database/sql/driver"
	"io"
	"math"

	"github.com/AlekSi/mysqlx/internal/mysqlx_resultset"
	"github.com/AlekSi/mysqlx/internal/mysqlx_sql"
)

// FIXME should it be thread-safe?
type rows struct {
	c       *conn
	columns []mysqlx_resultset.ColumnMetaData
	rows    chan *mysqlx_resultset.Row
	readErr error
}

func (r *rows) runReader() {
	defer close(r.rows)

	for {
		m, err := readMessage(r.c.transport)
		if err != nil {
			r.readErr = err
			return
		}

		switch m := m.(type) {
		case *mysqlx_resultset.Row:
			r.rows <- m
		case *mysqlx_resultset.FetchDone:
			continue
		case *mysqlx_sql.StmtExecuteOk:
			return
		default:
			bugf("unhandled message %T", m)
		}
	}
}

func (r *rows) Columns() []string {
	res := make([]string, len(r.columns))
	for i, c := range r.columns {
		res[i] = string(c.Name)
	}
	return res
}

func (r *rows) ColumnTypeDatabaseTypeName(index int) string {
	return r.columns[index].Type.String()
}

func (r *rows) ColumnTypeLength(index int) (length int64, ok bool) {
	c := r.columns[index]
	length = int64(c.GetLength())
	if c.Type.String() == "BYTES" && length == math.MaxUint32 {
		length = math.MaxInt64
	}
	ok = true
	return
}

func (r *rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	nullable = (*r.columns[index].Flags & 0x0010) != 0
	ok = true
	return
}

func (r *rows) Close() error {
	// drain messages
	for range r.rows {
	}
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	row, ok := <-r.rows
	if !ok {
		if r.readErr != nil {
			return r.readErr
		}
		return io.EOF
	}

	var err error
	for i, value := range row.Field {
		d, e := unmarshalValue(value, &r.columns[i])
		dest[i] = d
		if e != nil {
			err = e
		}
	}
	return err
}

// check interfaces
var (
	_ driver.Rows                           = (*rows)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*rows)(nil)
	_ driver.RowsColumnTypeLength           = (*rows)(nil)
	_ driver.RowsColumnTypeNullable         = (*rows)(nil)

	// TODO
	// _ driver.RowsColumnTypePrecisionScale = (*rows)(nil)
	// _ driver.RowsColumnTypeScanType       = (*rows)(nil)
	// _ driver.RowsNextResultSet            = (*rows)(nil)
)
