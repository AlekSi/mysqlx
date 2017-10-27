package mysqlx

import (
	"context"
	"database/sql/driver"
	"io"
	"math"

	"github.com/AlekSi/mysqlx/internal/mysqlx_notice"
	"github.com/AlekSi/mysqlx/internal/mysqlx_resultset"
	"github.com/AlekSi/mysqlx/internal/mysqlx_sql"
)

// rows is an iterator over an executed query's results.
type rows struct {
	c       *conn
	columns []mysqlx_resultset.ColumnMetaData
	rows    chan *mysqlx_resultset.Row
	readErr error
}

// runReader reads rows and sends then to channel until all rows are read or error is encountered.
func (r *rows) runReader(ctx context.Context) {
	defer close(r.rows)

	for {
		m, err := r.c.readMessage(ctx)
		if err != nil {
			r.readErr = err
			return
		}

		switch m := m.(type) {
		case *mysqlx_resultset.Row:
			r.rows <- m
		case *mysqlx_resultset.FetchDone:
			continue
		case *mysqlx_notice.SessionStateChanged:
			switch m.GetParam() {
			case mysqlx_notice.SessionStateChanged_ROWS_AFFECTED:
				continue
			default:
				r.readErr = bugf("rows.runReader: unhandled session state change %v", m)
				return
			}
		case *mysqlx_sql.StmtExecuteOk:
			return
		default:
			r.readErr = bugf("rows.runReader: unhandled message %T", m)
			return
		}
	}
}

// Columns returns the names of the columns.
// The number of columns of the result is inferred from the length of the slice.
// If a particular column name isn't known, an empty string should be returned for that entry.
func (r *rows) Columns() []string {
	res := make([]string, len(r.columns))
	for i, c := range r.columns {
		res[i] = string(c.Name)
	}
	return res
}

// Close closes the rows iterator.
func (r *rows) Close() error {
	// TODO limit a number of messages to drain there? and close connection?

	// drain messages until r.rows is closed in runReader
	for range r.rows {
	}

	// FIXME should we return r.readErr instead of nil?
	return nil
}

// Next is called to populate the next row of data into the provided slice.
// The provided slice will be the same size as the Columns() are wide.
// Next should return io.EOF when there are no more rows.
func (r *rows) Next(dest []driver.Value) error {
	row, ok := <-r.rows
	if !ok {
		if r.readErr != nil {
			return r.readErr
		}
		return io.EOF
	}

	// unmarshal all values, return first encountered error
	var err error
	for i, value := range row.Field {
		d, e := unmarshalValue(value, &r.columns[i])
		dest[i] = d
		if err == nil {
			err = e
		}
	}
	return err
}

// RowsColumnTypeDatabaseTypeName returns the database system type name without the length.
// Type names should be uppercase. Returned types:
// SINT, UINT, DOUBLE, FLOAT, BYTES, TIME, DATETIME, SET, ENUM, BIT, DECIMAL.
func (r *rows) ColumnTypeDatabaseTypeName(index int) string {
	return r.columns[index].Type.String()
}

// RowsColumnTypeLength return the length of the column type if the column is a variable length type.
// If the column is not a variable length type ok should return false.
// If length is not limited other than system limits, it should return math.MaxInt64.
// The following are examples of returned values for various types:
// TODO add examples
func (r *rows) ColumnTypeLength(index int) (length int64, ok bool) {
	c := r.columns[index]
	length = int64(c.GetLength())
	if c.Type.String() == "BYTES" && length == math.MaxUint32 {
		length = math.MaxInt64
	}
	ok = true
	return
}

// RowsColumnTypeNullable returns true if it is known the column may be null,
// or false if the column is known to be not nullable.
// If the column nullability is unknown, ok should be false.
func (r *rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	nullable = (r.columns[index].GetFlags() & 0x0010) != 0
	ok = true
	return
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
