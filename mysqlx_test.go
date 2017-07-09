package mysqlx

import (
	"database/sql"
	"math"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type City struct {
	ID       int
	Name     string
	District string
	Info     string
}

func (c *City) Values() []interface{} {
	return []interface{}{
		c.ID,
		c.Name,
		c.District,
		c.Info,
	}
}

func (c *City) Pointers() []interface{} {
	return []interface{}{
		&c.ID,
		&c.Name,
		&c.District,
		&c.Info,
	}
}

type ColumnType struct {
	Name             string
	DatabaseTypeName string
	Length           int64
	// ScanType         reflect.Type
	// TODO more checks
}

type MySQLXSuite struct {
	suite.Suite
	db *sql.DB
}

func TestMySQLX(t *testing.T) {
	suite.Run(t, new(MySQLXSuite))
}

func (s *MySQLXSuite) SetupTest() {
	testTraceF = s.T().Logf

	ds := os.Getenv("MYSQLX_TEST_DATASOURCE")
	s.Require().NotEmpty(ds, "Please set environment variable MYSQLX_TEST_DATASOURCE.")
	u, err := url.Parse(ds)
	s.Require().NoError(err)
	u.Path = "world_x"

	s.db, err = sql.Open("mysqlx", u.String())
	s.Require().NoError(err)
	s.Require().NoError(s.db.Ping())
}

func (s *MySQLXSuite) TearDownTest() {
	s.NoError(s.db.Close())
}

func (s *MySQLXSuite) TestQueryTable() {
	rows, err := s.db.Query("SELECT ID, Name, District, Info FROM city WHERE CountryCode = ? ORDER BY ID LIMIT 3", "RUS")
	s.Require().NoError(err)

	columns, err := rows.Columns()
	s.NoError(err)
	s.Equal([]string{"ID", "Name", "District", "Info"}, columns)

	types, err := rows.ColumnTypes()
	s.NoError(err)
	s.Require().Len(types, 4)
	for i, expected := range []ColumnType{
		// TODO convert internal X Protocol types to MySQL types (?)
		{"ID", "SINT", 11},        // FIXME should length of int(11) really be 11?
		{"Name", "BYTES", 35 * 3}, // CHAR(35) (inlike VARCHAR) stores 3 bytes per utf8 rune
		{"District", "BYTES", 20 * 3},
		{"Info", "BYTES", math.MaxInt64},
	} {
		s.Equal(expected.Name, types[i].Name(), "type %+v", types[i])
		s.Equal(expected.DatabaseTypeName, types[i].DatabaseTypeName(), "type %+v", types[i])
		l, ok := types[i].Length()
		if !ok {
			l = -1
		}
		s.Equal(expected.Length, l, "type %+v", types[i])
		// TODO more checks
	}

	for _, expected := range []City{
		{3580, "Moscow", "Moscow (City)", `{"Population": 8389200}`},
		{3581, "St Petersburg", "Pietari", `{"Population": 4694000}`},
		{3582, "Novosibirsk", "Novosibirsk", `{"Population": 1398800}`},
	} {
		s.True(rows.Next())
		var actual City
		s.NoError(rows.Scan(actual.Pointers()...))
		s.Equal(expected, actual)
	}

	s.False(rows.Next())
	s.NoError(rows.Err())
	s.NoError(rows.Close())
}

func (s *MySQLXSuite) TestQueryData() {
	fullDate := time.Date(2017, 7, 1, 12, 34, 56, 123456789, time.UTC)
	for _, q := range []struct {
		query    string
		arg      []interface{}
		expected interface{}
	}{
		// untyped NULL
		{`SELECT NULL`, nil, nil},
		{`SELECT ?`, []interface{}{nil}, nil},

		// CHAR
		{`SELECT 'foo'`, nil, "foo"},
		{`SELECT ?`, []interface{}{"foo"}, "foo"},
		{`SELECT ''`, nil, ""},
		{`SELECT ?`, []interface{}{""}, ""},
		{`SELECT CAST(NULL AS CHAR)`, nil, nil},
		{`SELECT CAST(? AS CHAR)`, []interface{}{nil}, nil},

		// SIGNED
		{`SELECT -42`, nil, int64(-42)},
		{`SELECT ?`, []interface{}{-42}, int64(-42)},
		{`SELECT -0`, nil, int64(0)},
		{`SELECT ?`, []interface{}{-0}, int64(0)},
		{`SELECT CAST(NULL AS SIGNED)`, nil, nil},
		{`SELECT CAST(? AS SIGNED)`, []interface{}{nil}, nil},

		// UNSIGNED
		{`SELECT CAST(42 AS UNSIGNED)`, nil, int64(42)},
		{`SELECT CAST(? AS UNSIGNED)`, []interface{}{42}, int64(42)},
		{`SELECT CAST(0 AS UNSIGNED)`, nil, int64(0)},
		{`SELECT CAST(? AS UNSIGNED)`, []interface{}{0}, int64(0)},
		{`SELECT CAST(NULL AS SIGNED)`, nil, nil},
		{`SELECT CAST(? AS SIGNED)`, []interface{}{nil}, nil},

		// DATE
		{`SELECT CAST('2017-07-01 12:34:56.123456789' AS DATE)`, nil, time.Date(2017, 7, 1, 0, 0, 0, 0, time.UTC)},
		{`SELECT CAST(? AS DATE)`, []interface{}{fullDate}, time.Date(2017, 7, 1, 0, 0, 0, 0, time.UTC)},
		{`SELECT CAST(NULL AS DATE)`, nil, nil},
		{`SELECT CAST(? AS DATE)`, []interface{}{nil}, nil},

		// DATETIME
		{`SELECT CAST('2017-07-01 12:34:56.123456789' AS DATETIME)`, nil, time.Date(2017, 7, 1, 12, 34, 56, 0, time.UTC)},
		{`SELECT CAST(? AS DATETIME)`, []interface{}{fullDate}, time.Date(2017, 7, 1, 12, 34, 56, 0, time.UTC)},
		{`SELECT CAST(NULL AS DATETIME)`, nil, nil},
		{`SELECT CAST(? AS DATETIME)`, []interface{}{nil}, nil},

		// DATETIME(6)
		{`SELECT CAST('2017-07-01 12:34:56.123456789' AS DATETIME(6))`, nil, time.Date(2017, 7, 1, 12, 34, 56, 123457000, time.UTC)},
		{`SELECT CAST(? AS DATETIME(6))`, []interface{}{fullDate}, time.Date(2017, 7, 1, 12, 34, 56, 123457000, time.UTC)},
		{`SELECT CAST(NULL AS DATETIME(6))`, nil, nil},
		{`SELECT CAST(? AS DATETIME(6))`, []interface{}{nil}, nil},
	} {
		// test QueryRow
		var actual interface{} = "NOT SET"
		s.Require().NoError(s.db.QueryRow(q.query, q.arg...).Scan(&actual))
		s.Equal(q.expected, actual)

		// test Query, read all rows
		rows, err := s.db.Query(q.query, q.arg...)
		s.Require().NoError(err)
		types, err := rows.ColumnTypes()
		s.NoError(err)
		s.Require().Len(types, 1)
		s.T().Log(types[0])

		s.True(rows.Next())
		actual = "NOT SET"
		s.NoError(rows.Scan(&actual))
		s.Equal(q.expected, actual)
		s.False(rows.Next())
		s.NoError(rows.Err())
		s.NoError(rows.Close())
	}
}

func (s *MySQLXSuite) TestQueryEmpty() {
	_, err := s.db.Exec("CREATE TEMPORARY TABLE TestQueryEmpty (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	s.Require().NoError(err)

	var actual interface{} = "NOT SET"
	s.Equal(sql.ErrNoRows, s.db.QueryRow("SELECT * FROM TestQueryEmpty").Scan(&actual))
	s.Equal("NOT SET", actual)
}

func (s *MySQLXSuite) TestQueryExec() {
	// test QueryRow
	var actual interface{} = "NOT SET"
	s.Equal(sql.ErrNoRows, s.db.QueryRow(`CREATE TEMPORARY TABLE TestQueryExec1 (id int)`).Scan(&actual))
	s.Equal("NOT SET", actual)

	// test Query, read all rows
	rows, err := s.db.Query(`CREATE TEMPORARY TABLE TestQueryExec2 (id int)`)
	s.Require().NoError(err)
	types, err := rows.ColumnTypes()
	s.NoError(err)
	s.Len(types, 0)
	s.False(rows.Next())
	s.NoError(rows.Err())
	s.NoError(rows.Close())
}

func (s *MySQLXSuite) TestExec() {
	res, err := s.db.Exec("CREATE TEMPORARY TABLE TestExec (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	s.Require().NoError(err)
	id, err := res.LastInsertId()
	s.EqualError(err, "no LastInsertId available")
	s.Equal(int64(0), id)
	ra, err := res.RowsAffected()
	s.NoError(err)
	s.Equal(int64(0), ra)

	res, err = s.db.Exec("INSERT INTO TestExec VALUES (1), (2)")
	s.Require().NoError(err)
	id, err = res.LastInsertId()
	s.NoError(err)
	s.Equal(int64(2), id)
	ra, err = res.RowsAffected()
	s.NoError(err)
	s.Equal(int64(2), ra)

	res, err = s.db.Exec("UPDATE TestExec SET id = ? WHERE id = ?", 3, 2)
	s.Require().NoError(err)
	id, err = res.LastInsertId()
	s.EqualError(err, "no LastInsertId available")
	s.Equal(int64(0), id)
	ra, err = res.RowsAffected()
	s.NoError(err)
	s.Equal(int64(1), ra)
}

func (s *MySQLXSuite) TestExecQuery() {
	res, err := s.db.Exec("SELECT 1")
	s.NoError(err)
	id, err := res.LastInsertId()
	s.EqualError(err, "no LastInsertId available")
	s.Equal(int64(0), id)
	ra, err := res.RowsAffected()
	s.NoError(err)
	s.Equal(int64(0), ra)
}

func (s *MySQLXSuite) TestBeginCommit() {
	_, err := s.db.Exec("CREATE TEMPORARY TABLE TestBeginRollback (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	s.Require().NoError(err)

	tx, err := s.db.Begin()
	s.Require().NoError(err)

	_, err = tx.Exec("INSERT INTO TestBeginRollback VALUES (1)")
	s.NoError(err)

	s.NoError(tx.Commit())
	s.Equal(sql.ErrTxDone, tx.Commit())
	s.Equal(sql.ErrTxDone, tx.Rollback())

	var count int
	s.NoError(s.db.QueryRow("SELECT COUNT(*) FROM TestBeginRollback").Scan(&count))
	s.Equal(1, count)
}

func (s *MySQLXSuite) TestBeginRollback() {
	_, err := s.db.Exec("CREATE TEMPORARY TABLE TestBeginRollback (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	s.Require().NoError(err)

	tx, err := s.db.Begin()
	s.Require().NoError(err)

	_, err = tx.Exec("INSERT INTO TestBeginRollback VALUES (1)")
	s.NoError(err)

	s.NoError(tx.Rollback())
	s.Equal(sql.ErrTxDone, tx.Rollback())
	s.Equal(sql.ErrTxDone, tx.Commit())

	var count int
	s.NoError(s.db.QueryRow("SELECT COUNT(*) FROM TestBeginRollback").Scan(&count))
	s.Equal(0, count)
}

func TestNoDatabase(t *testing.T) {
	testTraceF = t.Logf

	ds := os.Getenv("MYSQLX_TEST_DATASOURCE")
	require.NotEmpty(t, ds, "Please set environment variable MYSQLX_TEST_DATASOURCE.")
	u, err := url.Parse(ds)
	require.NoError(t, err)

	db, err := sql.Open("mysqlx", u.String())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, db.Close())
	}()

	require.NoError(t, db.Ping())

	var s string
	require.NoError(t, db.QueryRow("SELECT VERSION()").Scan(&s))
	t.Log(s)

	err = db.QueryRow("SELECT Name FROM city LIMIT 1").Scan(&s)
	assert.Equal(t, &Error{Severity: SeverityError, Code: 1046, SQLState: "3D000", Msg: "No database selected"}, err)
	assert.Equal(t, "ERROR 1046 (3D000): No database selected", err.Error())

	res, err := db.Exec("UPDATE city SET Name = ?", "Moscow")
	assert.Nil(t, res)
	assert.Equal(t, &Error{Severity: SeverityError, Code: 1046, SQLState: "3D000", Msg: "No database selected"}, err)
	assert.Equal(t, "ERROR 1046 (3D000): No database selected", err.Error())
}
