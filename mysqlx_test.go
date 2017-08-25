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

type CountryLanguage struct {
	CountryCode string
	Language    string
	IsOfficial  bool
	Percentage  float32
}

func (cl *CountryLanguage) Values() []interface{} {
	return []interface{}{
		cl.CountryCode,
		cl.Language,
		cl.IsOfficial,
		cl.Percentage,
	}
}

func (cl *CountryLanguage) Pointers() []interface{} {
	return []interface{}{
		&cl.CountryCode,
		&cl.Language,
		&cl.IsOfficial,
		&cl.Percentage,
	}
}

type ColumnType struct {
	Name             string
	DatabaseTypeName string
	Length           int64
	// ScanType         reflect.Type
	// TODO more checks
}

func openDB(t *testing.T, database string) *sql.DB {
	t.Helper()
	setTestTracef(t.Name(), t.Logf)

	ds := os.Getenv("MYSQLX_TEST_DATASOURCE")
	require.NotEmpty(t, ds, "Please set environment variable MYSQLX_TEST_DATASOURCE.")
	u, err := url.Parse(ds)
	require.NoError(t, err)
	u.Path = database
	q := u.Query()
	q.Set("_trace", t.Name())
	u.RawQuery = q.Encode()

	db, err := sql.Open("mysqlx", u.String())
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func closeDB(t *testing.T, db *sql.DB) {
	assert.NoError(t, db.Close())
}

func TestQueryTableCity(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	rows, err := db.Query("SELECT ID, Name, District, Info FROM city WHERE CountryCode = ? ORDER BY ID LIMIT 3", "RUS")
	require.NoError(t, err)

	columns, err := rows.Columns()
	assert.NoError(t, err)
	assert.Equal(t, []string{"ID", "Name", "District", "Info"}, columns)

	types, err := rows.ColumnTypes()
	assert.NoError(t, err)
	require.Len(t, types, 4)
	for i, expected := range []ColumnType{
		// TODO convert internal X Protocol types to MySQL types (?)
		{"ID", "SINT", 11},        // FIXME should length of int(11) really be 11?
		{"Name", "BYTES", 35 * 3}, // CHAR(35) (inlike VARCHAR) stores 3 bytes per utf8 rune
		{"District", "BYTES", 20 * 3},
		{"Info", "BYTES", math.MaxInt64},
	} {
		assert.Equal(t, expected.Name, types[i].Name(), "type %+v", types[i])
		assert.Equal(t, expected.DatabaseTypeName, types[i].DatabaseTypeName(), "type %+v", types[i])
		l, ok := types[i].Length()
		if !ok {
			l = -1
		}
		assert.Equal(t, expected.Length, l, "type %+v", types[i])
		// TODO more checks
	}

	for _, expected := range []City{
		{3580, "Moscow", "Moscow (City)", `{"Population": 8389200}`},
		{3581, "St Petersburg", "Pietari", `{"Population": 4694000}`},
		{3582, "Novosibirsk", "Novosibirsk", `{"Population": 1398800}`},
	} {
		assert.True(t, rows.Next())
		var actual City
		assert.NoError(t, rows.Scan(actual.Pointers()...))
		assert.Equal(t, expected, actual)
	}

	assert.False(t, rows.Next())
	assert.NoError(t, rows.Err())
	assert.NoError(t, rows.Close())
}

func TestQueryTableCountryLanguage(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	rows, err := db.Query("SELECT CountryCode, Language, IsOfficial, Percentage FROM countrylanguage WHERE CountryCode = ? ORDER BY Percentage DESC LIMIT 3", "RUS")
	require.NoError(t, err)

	columns, err := rows.Columns()
	assert.NoError(t, err)
	assert.Equal(t, []string{"CountryCode", "Language", "IsOfficial", "Percentage"}, columns)

	types, err := rows.ColumnTypes()
	assert.NoError(t, err)
	require.Len(t, types, 4)
	for i, expected := range []ColumnType{
		// TODO convert internal X Protocol types to MySQL types (?)
		{"CountryCode", "BYTES", 3 * 3}, // CHAR(3) (inlike VARCHAR) stores 3 bytes per utf8 rune
		{"Language", "BYTES", 30 * 3},
		{"IsOfficial", "ENUM", 1 * 3},
		{"Percentage", "FLOAT", 4},
	} {
		assert.Equal(t, expected.Name, types[i].Name(), "type %+v", types[i])
		assert.Equal(t, expected.DatabaseTypeName, types[i].DatabaseTypeName(), "type %+v", types[i])
		l, ok := types[i].Length()
		if !ok {
			l = -1
		}
		assert.Equal(t, expected.Length, l, "type %+v", types[i])
		// TODO more checks
	}

	for _, expected := range []CountryLanguage{
		{"RUS", "Russian", true, 86.6},
		{"RUS", "Tatar", false, 3.2},
		{"RUS", "Ukrainian", false, 1.3},
	} {
		assert.True(t, rows.Next())
		var actual CountryLanguage
		assert.NoError(t, rows.Scan(actual.Pointers()...))
		assert.Equal(t, expected, actual)
	}

	assert.False(t, rows.Next())
	assert.NoError(t, rows.Err())
	assert.NoError(t, rows.Close())
}

func TestQueryData(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

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
		{`SELECT CAST(NULL AS UNSIGNED)`, nil, nil},
		{`SELECT CAST(? AS UNSIGNED)`, []interface{}{nil}, nil},

		// floats are returned as strings
		{`SELECT 12.3401`, nil, "12.3401"},
		{`SELECT ?`, []interface{}{12.3401}, "12.3401"},

		// DECIMAL
		// values from datatypes_test.go
		{`SELECT CAST(12.3401 AS DECIMAL(6,4))`, nil, "12.3401"},
		{`SELECT CAST(-12.3401 AS DECIMAL(6,4))`, nil, "-12.3401"},
		{`SELECT CAST(12.3401 AS DECIMAL(6,3))`, nil, "12.340"},
		{`SELECT CAST(-12.3401 AS DECIMAL(6,3))`, nil, "-12.340"},
		{`SELECT CAST(12.3401 AS DECIMAL(6))`, nil, "12"},
		{`SELECT CAST(-12.3401 AS DECIMAL(6))`, nil, "-12"},
		{`SELECT CAST(12.3401 AS DECIMAL(1))`, nil, "9"},       // Warning (code 1264): Out of range value
		{`SELECT CAST(-12.3401 AS DECIMAL(1))`, nil, "-9"},     // Warning (code 1264): Out of range value
		{`SELECT CAST(12.3401 AS DECIMAL(1,1))`, nil, "0.9"},   // Warning (code 1264): Out of range value
		{`SELECT CAST(-12.3401 AS DECIMAL(1,1))`, nil, "-0.9"}, // Warning (code 1264): Out of range value

		{`SELECT CAST(? AS DECIMAL(6,4))`, []interface{}{"-12.3401"}, "-12.3401"},
		{`SELECT CAST(0 AS DECIMAL(6,4))`, nil, "0.0000"},
		{`SELECT CAST(? AS DECIMAL(6,4))`, []interface{}{"0"}, "0.0000"},
		{`SELECT CAST(NULL AS DECIMAL(6,4))`, nil, nil},
		{`SELECT CAST(? AS DECIMAL(6,4))`, []interface{}{nil}, nil},

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
		require.NoError(t, db.QueryRow(q.query, q.arg...).Scan(&actual))
		assert.Equal(t, q.expected, actual)

		// test Query, read all rows
		rows, err := db.Query(q.query, q.arg...)
		require.NoError(t, err)
		types, err := rows.ColumnTypes()
		assert.NoError(t, err)
		require.Len(t, types, 1)
		t.Log(types[0])

		assert.True(t, rows.Next())
		actual = "NOT SET"
		assert.NoError(t, rows.Scan(&actual))
		assert.Equal(t, q.expected, actual)
		assert.False(t, rows.Next())
		assert.NoError(t, rows.Err())
		assert.NoError(t, rows.Close())

		stmt, err := db.Prepare(q.query)
		require.NoError(t, err)

		// test Prepare + QueryRow
		actual = "NOT SET"
		require.NoError(t, stmt.QueryRow(q.arg...).Scan(&actual))
		assert.Equal(t, q.expected, actual)

		// test Prepare + Query, read all rows
		rows, err = stmt.Query(q.arg...)
		require.NoError(t, err)
		types, err = rows.ColumnTypes()
		assert.NoError(t, err)
		require.Len(t, types, 1)
		t.Log(types[0])

		assert.True(t, rows.Next())
		actual = "NOT SET"
		assert.NoError(t, rows.Scan(&actual))
		assert.Equal(t, q.expected, actual)
		assert.False(t, rows.Next())
		assert.NoError(t, rows.Err())
		assert.NoError(t, rows.Close())

		assert.NoError(t, stmt.Close())
	}
}

func TestQueryEmpty(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	_, err := db.Exec("CREATE TEMPORARY TABLE TestQueryEmpty (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	require.NoError(t, err)

	var actual interface{} = "NOT SET"
	assert.Equal(t, sql.ErrNoRows, db.QueryRow("SELECT * FROM TestQueryEmpty").Scan(&actual))
	assert.Equal(t, "NOT SET", actual)
}

func TestQueryExec(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	// test QueryRow
	var actual interface{} = "NOT SET"
	assert.Equal(t, sql.ErrNoRows, db.QueryRow(`CREATE TEMPORARY TABLE TestQueryExec1 (id int)`).Scan(&actual))
	assert.Equal(t, "NOT SET", actual)

	// test Query, read all rows
	rows, err := db.Query(`CREATE TEMPORARY TABLE TestQueryExec2 (id int)`)
	require.NoError(t, err)
	types, err := rows.ColumnTypes()
	assert.NoError(t, err)
	assert.Len(t, types, 0)
	assert.False(t, rows.Next())
	assert.NoError(t, rows.Err())
	assert.NoError(t, rows.Close())
}

func TestQueryCloseEarly(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	// read 0 rows
	rows, err := db.Query("SELECT ID, Name, District, Info FROM city WHERE CountryCode = ? ORDER BY ID LIMIT 3", "RUS")
	require.NoError(t, err)
	assert.NoError(t, rows.Close())

	// read 1 row
	var city City
	rows, err = db.Query("SELECT ID, Name, District, Info FROM city WHERE CountryCode = ? ORDER BY ID LIMIT 3", "USA")
	require.NoError(t, err)
	assert.True(t, rows.Next())
	assert.NoError(t, rows.Scan(city.Pointers()...))
	assert.Equal(t, City{3793, "New York", "New York", `{"Population": 8008278}`}, city)
	assert.NoError(t, rows.Close())

	// read 2 rows
	rows, err = db.Query("SELECT ID, Name, District, Info FROM city WHERE CountryCode = ? ORDER BY ID LIMIT 3", "FRA")
	require.NoError(t, err)
	assert.True(t, rows.Next())
	assert.NoError(t, rows.Scan(city.Pointers()...))
	assert.Equal(t, City{2974, "Paris", "Île-de-France", `{"Population": 2125246}`}, city)
	assert.True(t, rows.Next())
	assert.NoError(t, rows.Scan(city.Pointers()...))
	assert.Equal(t, City{2975, "Marseille", "Provence-Alpes-Côte", `{"Population": 798430}`}, city)
	assert.NoError(t, rows.Close())
}

func TestExec(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	res, err := db.Exec("CREATE TEMPORARY TABLE TestExec (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	require.NoError(t, err)
	id, err := res.LastInsertId()
	assert.EqualError(t, err, "no LastInsertId available")
	assert.Equal(t, int64(0), id)
	ra, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), ra)

	res, err = db.Exec("INSERT INTO TestExec VALUES (1), (2)")
	require.NoError(t, err)
	id, err = res.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), id)
	ra, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), ra)

	res, err = db.Exec("UPDATE TestExec SET id = ? WHERE id = ?", 3, 1)
	require.NoError(t, err)
	id, err = res.LastInsertId()
	assert.EqualError(t, err, "no LastInsertId available")
	assert.Equal(t, int64(0), id)
	ra, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ra)

	stmt, err := db.Prepare("UPDATE TestExec SET id = ? WHERE id = ?")
	require.NoError(t, err)
	res, err = stmt.Exec(4, 2)
	require.NoError(t, err)
	id, err = res.LastInsertId()
	assert.EqualError(t, err, "no LastInsertId available")
	assert.Equal(t, int64(0), id)
	ra, err = res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ra)
	assert.NoError(t, stmt.Close())
}

func TestExecQuery(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	res, err := db.Exec("SELECT 1")
	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.EqualError(t, err, "no LastInsertId available")
	assert.Equal(t, int64(0), id)
	ra, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), ra)
}

func TestBeginCommit(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	_, err := db.Exec("CREATE TEMPORARY TABLE TestBeginCommit (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	require.NoError(t, err)

	tx, err := db.Begin()
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO TestBeginCommit VALUES (1)")
	assert.NoError(t, err)

	assert.NoError(t, tx.Commit())
	assert.Equal(t, sql.ErrTxDone, tx.Commit())
	assert.Equal(t, sql.ErrTxDone, tx.Rollback())

	var count int
	assert.NoError(t, db.QueryRow("SELECT COUNT(*) FROM TestBeginCommit").Scan(&count))
	assert.Equal(t, 1, count)
}

func TestBeginRollback(t *testing.T) {
	t.Parallel()
	db := openDB(t, "world_x")
	defer closeDB(t, db)

	_, err := db.Exec("CREATE TEMPORARY TABLE TestBeginRollback (id int AUTO_INCREMENT, PRIMARY KEY (id))")
	require.NoError(t, err)

	tx, err := db.Begin()
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO TestBeginRollback VALUES (1)")
	assert.NoError(t, err)

	assert.NoError(t, tx.Rollback())
	assert.Equal(t, sql.ErrTxDone, tx.Rollback())
	assert.Equal(t, sql.ErrTxDone, tx.Commit())

	var count int
	assert.NoError(t, db.QueryRow("SELECT COUNT(*) FROM TestBeginRollback").Scan(&count))
	assert.Equal(t, 0, count)
}

func TestNoDatabase(t *testing.T) {
	t.Parallel()
	db := openDB(t, "")
	defer closeDB(t, db)

	var s string
	require.NoError(t, db.QueryRow("SELECT VERSION()").Scan(&s))
	t.Log(s)

	err := db.QueryRow("SELECT Name FROM city LIMIT 1").Scan(&s)
	assert.Equal(t, &Error{Severity: SeverityError, Code: 1046, SQLState: "3D000", Msg: "No database selected"}, err)
	assert.EqualError(t, err, "ERROR 1046 (3D000): No database selected")

	res, err := db.Exec("UPDATE city SET Name = ?", "Moscow")
	assert.Nil(t, res)
	assert.Equal(t, &Error{Severity: SeverityError, Code: 1046, SQLState: "3D000", Msg: "No database selected"}, err)
	assert.EqualError(t, err, "ERROR 1046 (3D000): No database selected")
}
