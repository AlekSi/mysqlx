package mysqlx

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/AlekSi/mysqlx/internal/mysqlx"
	"github.com/AlekSi/mysqlx/internal/mysqlx_connection"
	"github.com/AlekSi/mysqlx/internal/mysqlx_notice"
	"github.com/AlekSi/mysqlx/internal/mysqlx_resultset"
	"github.com/AlekSi/mysqlx/internal/mysqlx_session"
	"github.com/AlekSi/mysqlx/internal/mysqlx_sql"
)

// TODO make this configurable?
// It should not be less then 1.
const rowsCap = 1

// https://github.com/mysql/mysql-server/blob/mysql-5.7.19/rapid/plugin/x/mysqlxtest_src/password_hasher.cc
// https://github.com/mysql/mysql-server/blob/mysql-5.7.19/rapid/plugin/x/mysqlxtest_src/mysql41_hash.cc
func scramble(data []byte, password string) []byte {
	hash1 := sha1.Sum([]byte(password))
	hash2 := sha1.Sum([]byte(hash1[:]))

	h := sha1.New()
	h.Write(data)
	h.Write(hash2[:])
	res := h.Sum(nil)

	for i := range res {
		res[i] ^= hash1[i]
	}
	return res[:]
}

func authData(database, username, password string, authData []byte) []byte {
	if len(authData) != 20 {
		return []byte(bugf("authData: expected authData to has 20 bytes, got %d", len(authData)).Error())
	}

	res := database + "\x00" + username + "\x00"
	if password == "" {
		return []byte(res)
	}

	res += fmt.Sprintf("*%02x", scramble(authData, password))
	return []byte(res)
}

// conn is a connection to a database.
// It is not used concurrently by multiple goroutines.
// conn is assumed to be stateful.
type conn struct {
	transport net.Conn
	tracef    TraceFunc

	closeOnce sync.Once
	closeErr  error
}

func newConn(transport net.Conn, traceF TraceFunc) *conn {
	traceF("+++ connection created: %s->%s", transport.LocalAddr(), transport.RemoteAddr())

	return &conn{
		transport: transport,
		tracef:    traceF,
	}
}

func setDefaults(u *url.URL) error {
	if u.Opaque != "" {
		return fmt.Errorf("invalid data source: %s", u.String())
	}

	if u.Scheme == "" {
		u.Scheme = "tcp"
	}
	host, port, _ := net.SplitHostPort(u.Host)
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "33060"
	}
	u.Host = net.JoinHostPort(host, port)

	return nil
}

func open(ctx context.Context, dataSource string, openParams *OpenParams) (*conn, error) {
	u, err := url.Parse(dataSource)
	if err != nil {
		return nil, err
	}
	if err = setDefaults(u); err != nil {
		return nil, err
	}

	// check and handle parameters, extract session variables
	params := u.Query()
	vars := make(map[string]string, len(params))
	traceF := openParams.Trace
	for k, vs := range params {
		if len(vs) != 1 {
			return nil, fmt.Errorf("%d values for parameter %q", len(vs), k)
		}
		v := vs[0]

		if !strings.HasPrefix(k, "_") {
			vars[k] = v
			continue
		}

		switch k {
		case "_trace":
			traceF = getTracef(v)
		case "_open_timeout":
			d, err := time.ParseDuration(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %s", k, err)
			}
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, d)
			defer cancel()
		default:
			return nil, fmt.Errorf("unexpected parameter %q", k)
		}
	}

	conn, err := openParams.Dialer.DialContext(ctx, u.Scheme, u.Host)
	if err != nil {
		return nil, err
	}
	c := newConn(conn, traceF)

	database := strings.TrimPrefix(u.Path, "/")
	var username, password string
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	if err = c.negotiate(ctx); err != nil {
		return nil, err
	}
	if err = c.auth(ctx, database, username, password); err != nil {
		return nil, err
	}

	// set session variables
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		nv := []driver.NamedValue{{
			Ordinal: 1,
			Value:   vars[k],
		}}
		if _, err = c.ExecContext(ctx, "SET SESSION "+k+" = ?", nv); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *conn) negotiate(ctx context.Context) error {
	if err := c.writeMessage(ctx, &mysqlx_connection.CapabilitiesGet{}); err != nil {
		return c.close(err)
	}
	m, err := c.readMessage(ctx)
	if err != nil {
		return c.close(err)
	}

	var mechs []string
	var mysql41Found bool
	for _, cap := range m.(*mysqlx_connection.Capabilities).Capabilities {
		if cap.GetName() == "authentication.mechanisms" {
			for _, value := range cap.Value.Array.Value {
				s := string(value.Scalar.VString.Value)
				if s == "MYSQL41" {
					mysql41Found = true
				}
				mechs = append(mechs, s)
			}
		}
	}

	if !mysql41Found {
		return fmt.Errorf("no MYSQL41 authentication mechanism: %v", mechs)
	}
	return nil
}

func (c *conn) auth(ctx context.Context, database, username, password string) error {
	// TODO use password

	mechName := "MYSQL41"
	if err := c.writeMessage(ctx, &mysqlx_session.AuthenticateStart{
		MechName: &mechName,
	}); err != nil {
		return c.close(err)
	}

	m, err := c.readMessage(ctx)
	if err != nil {
		return c.close(err)
	}
	cont := m.(*mysqlx_session.AuthenticateContinue)

	if err = c.writeMessage(ctx, &mysqlx_session.AuthenticateContinue{
		AuthData: authData(database, username, password, cont.AuthData),
	}); err != nil {
		return c.close(err)
	}

	if m, err = c.readMessage(ctx); err != nil {
		return c.close(err)
	}
	_ = m.(*mysqlx_notice.SessionStateChanged)

	if m, err = c.readMessage(ctx); err != nil {
		return c.close(err)
	}
	_ = m.(*mysqlx_session.AuthenticateOk)

	return nil
}

func (c *conn) close(err error) error {
	c.closeOnce.Do(func() {
		c.closeErr = err
		e := c.transport.Close()
		if c.closeErr == nil {
			c.closeErr = e
		}
	})

	c.tracef("--- connection closed: %s->%s", c.transport.LocalAddr(), c.transport.RemoteAddr())
	return c.closeErr
}

// Close invalidates and potentially stops any current prepared statements and transactions,
// marking this connection as no longer in use.
// Because the sql package maintains a free pool of connections and only calls Close when there's
// a surplus of idle connections, it shouldn't be necessary for drivers to do their own connection caching.
func (c *conn) Close() error {
	if err := c.writeMessage(context.TODO(), &mysqlx_connection.Close{}); err != nil {
		return c.close(err)
	}

	// read one next message, but do not check it is mysqlx.Ok
	if _, err := c.readMessage(context.TODO()); err != nil {
		return c.close(err)
	}

	return c.close(nil)
}

// Begin starts and returns a new transaction.
func (c *conn) Begin() (driver.Tx, error) {
	if _, err := c.Exec("BEGIN", nil); err != nil {
		return nil, err
	}
	return &tx{
		c: c,
	}, nil
}

// BeginTx starts and returns a new transaction.
// If the context is canceled by the user the sql package will
// call Tx.Rollback before discarding and closing the connection.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if sql.IsolationLevel(opts.Isolation) != sql.LevelDefault {
		return nil, bugf("conn.BeginTx: isolation level %d is not supported yet", opts.Isolation)
	}
	if opts.ReadOnly {
		return nil, bugf("conn.BeginTx: read-only transactions are not supported yet")
	}
	if _, err := c.ExecContext(ctx, "BEGIN", nil); err != nil {
		return nil, err
	}
	return &tx{
		c: c,
	}, nil
}

// Prepare returns a prepared statement, bound to this connection.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return &stmt{
		c:     c,
		query: query,
	}, nil
}

// PrepareContext returns a prepared statement, bound to this connection.
// context is for the preparation of the statement (and so ignored).
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return &stmt{
		c:     c,
		query: query,
	}, nil
}

// Exec executes a query that doesn't return rows, such as an INSERT or UPDATE.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	nv := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		nv[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   arg,
		}
	}
	return c.ExecContext(context.Background(), query, nv)
}

// ExecContext executes a query that doesn't return rows, such as an INSERT or UPDATE.
// It honors the context timeout and return when the context is canceled.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	stmt := &mysqlx_sql.StmtExecute{
		Stmt: []byte(query),
	}
	for i, nv := range args {
		if nv.Name != "" {
			return nil, bugf("conn.ExecContext: %q - named values are not supported yet", nv.Name)
		}
		if nv.Ordinal != i+1 {
			return nil, bugf("conn.ExecContext: out-of-order values are not supported yet")
		}
		a, err := marshalValue(nv.Value)
		if err != nil {
			return nil, err
		}
		stmt.Args = append(stmt.Args, a)
	}

	if err := c.writeMessage(ctx, stmt); err != nil {
		return nil, c.close(err)
	}

	var result driver.Result = driver.ResultNoRows
	for {
		m, err := c.readMessage(ctx)
		if err != nil {
			return nil, c.close(err)
		}

		switch m := m.(type) {
		case *mysqlx.Error:
			severity := Severity(m.GetSeverity())

			// TODO close connection if severity is FATAL?

			return nil, &Error{
				Severity: severity,
				Code:     m.GetCode(),
				SQLState: m.GetSqlState(),
				Msg:      m.GetMsg(),
			}

		case *mysqlx_resultset.ColumnMetaData:
			continue

		// query with rows
		case *mysqlx_resultset.Row:
			continue

		// query without rows
		case *mysqlx_resultset.FetchDone:
			continue
		case *mysqlx_notice.SessionStateChanged:
			switch m.GetParam() {
			case mysqlx_notice.SessionStateChanged_GENERATED_INSERT_ID:
				ra, _ := result.RowsAffected()
				result = execResult{
					lastInsertId: int64(m.GetValue().GetVUnsignedInt()),
					rowsAffected: ra,
				}
			case mysqlx_notice.SessionStateChanged_ROWS_AFFECTED:
				if result == driver.ResultNoRows {
					result = driver.RowsAffected(m.GetValue().GetVUnsignedInt())
				}
			case mysqlx_notice.SessionStateChanged_PRODUCED_MESSAGE:
				// TODO log it?
				continue
			default:
				return nil, bugf("conn.Exec: unhandled session state change %v", m)
			}
		case *mysqlx_sql.StmtExecuteOk:
			return result, nil

		default:
			return nil, bugf("conn.Exec: unhandled type %T", m)
		}
	}
}

// Query executes a query that may return rows, such as a SELECT.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	nv := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		nv[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   arg,
		}
	}
	return c.QueryContext(context.Background(), query, nv)
}

// QueryContext executes a query that may return rows, such as a SELECT.
// It honors the context timeout and return when the context is canceled.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	stmt := &mysqlx_sql.StmtExecute{
		Stmt: []byte(query),
	}
	for i, nv := range args {
		if nv.Name != "" {
			return nil, bugf("conn.ExecContext: %q - named values are not supported yet", nv.Name)
		}
		if nv.Ordinal != i+1 {
			return nil, bugf("conn.ExecContext: out-of-order values are not supported yet")
		}
		a, err := marshalValue(nv.Value)
		if err != nil {
			return nil, err
		}
		stmt.Args = append(stmt.Args, a)
	}

	if err := c.writeMessage(ctx, stmt); err != nil {
		return nil, c.close(err)
	}

	rows := rows{
		c:       c,
		columns: make([]mysqlx_resultset.ColumnMetaData, 0, 1),
		rows:    make(chan *mysqlx_resultset.Row, rowsCap),
	}
	for {
		m, err := c.readMessage(ctx)
		if err != nil {
			return nil, c.close(err)
		}

		switch m := m.(type) {
		case *mysqlx.Error:
			severity := Severity(m.GetSeverity())

			// TODO close connection if severity is FATAL?

			return nil, &Error{
				Severity: severity,
				Code:     m.GetCode(),
				SQLState: m.GetSqlState(),
				Msg:      m.GetMsg(),
			}

		case *mysqlx_resultset.ColumnMetaData:
			rows.columns = append(rows.columns, *m)

		// query with rows
		case *mysqlx_resultset.Row:
			rows.rows <- m
			go rows.runReader(ctx)
			return &rows, nil

		// query without rows
		case *mysqlx_resultset.FetchDone:
			continue
		case *mysqlx_notice.SessionStateChanged:
			switch m.GetParam() {
			case mysqlx_notice.SessionStateChanged_ROWS_AFFECTED:
				continue
			default:
				return nil, bugf("conn.Query: unhandled session state change %v", m)
			}
		case *mysqlx_sql.StmtExecuteOk:
			close(rows.rows)
			return &rows, nil

		default:
			return nil, bugf("conn.Query: unhandled type %T", m)
		}
	}
}

func (c *conn) Ping(ctx context.Context) error {
	if _, err := c.ExecContext(ctx, "SELECT 'ping'", nil); err != nil {
		return driver.ErrBadConn
	}
	return nil
}

func (c *conn) writeMessage(ctx context.Context, m proto.Message) error {
	deadline, _ := ctx.Deadline()
	if err := c.transport.SetWriteDeadline(deadline); err != nil {
		return err
	}

	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	var t mysqlx.ClientMessages_Type
	switch m.(type) {
	case *mysqlx_connection.CapabilitiesGet:
		t = mysqlx.ClientMessages_CON_CAPABILITIES_GET
	case *mysqlx_connection.Close:
		t = mysqlx.ClientMessages_CON_CLOSE

	case *mysqlx_session.AuthenticateStart:
		t = mysqlx.ClientMessages_SESS_AUTHENTICATE_START
	case *mysqlx_session.AuthenticateContinue:
		t = mysqlx.ClientMessages_SESS_AUTHENTICATE_CONTINUE

	case *mysqlx_sql.StmtExecute:
		t = mysqlx.ClientMessages_SQL_STMT_EXECUTE

	default:
		return bugf("conn.writeMessage: unhandled client message: %T %#v", m, m)
	}

	c.tracef(">>> %T %v", m, m)

	var head [5]byte
	binary.LittleEndian.PutUint32(head[:], uint32(len(b))+1)
	head[4] = byte(t)
	_, err = (&net.Buffers{head[:], b}).WriteTo(c.transport)
	return err
}

func (c *conn) readMessage(ctx context.Context) (proto.Message, error) {
	deadline, _ := ctx.Deadline()
	if err := c.transport.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	var head [5]byte
	if _, err := io.ReadFull(c.transport, head[:]); err != nil {
		return nil, err
	}
	l := binary.LittleEndian.Uint32(head[:])
	t := mysqlx.ServerMessages_Type(head[4])

	var m proto.Message
	switch t {
	case mysqlx.ServerMessages_OK:
		m = new(mysqlx.Ok)
	case mysqlx.ServerMessages_ERROR:
		m = new(mysqlx.Error)

	case mysqlx.ServerMessages_CONN_CAPABILITIES:
		m = new(mysqlx_connection.Capabilities)

	case mysqlx.ServerMessages_SESS_AUTHENTICATE_CONTINUE:
		m = new(mysqlx_session.AuthenticateContinue)
	case mysqlx.ServerMessages_SESS_AUTHENTICATE_OK:
		m = new(mysqlx_session.AuthenticateOk)

	case mysqlx.ServerMessages_NOTICE:
		m = new(mysqlx_notice.Frame)

	case mysqlx.ServerMessages_RESULTSET_COLUMN_META_DATA:
		m = new(mysqlx_resultset.ColumnMetaData)
	case mysqlx.ServerMessages_RESULTSET_ROW:
		m = new(mysqlx_resultset.Row)
	case mysqlx.ServerMessages_RESULTSET_FETCH_DONE:
		// TODO short circuit there
		m = new(mysqlx_resultset.FetchDone)
	// case mysqlx.ServerMessages_RESULTSET_FETCH_SUSPENDED:
	// 	// FIXME what's there?
	// case mysqlx.ServerMessages_RESULTSET_FETCH_DONE_MORE_RESULTSETS:
	// 	m = new(mysqlx_resultset.FetchDoneMoreResultsets)
	// case mysqlx.ServerMessages_RESULTSET_FETCH_DONE_MORE_OUT_PARAMS:
	// 	m = new(mysqlx_resultset.FetchDoneMoreOutParams)

	case mysqlx.ServerMessages_SQL_STMT_EXECUTE_OK:
		// TODO short circuit there
		m = new(mysqlx_sql.StmtExecuteOk)

	default:
		return nil, bugf("conn.readMessage: unhandled type of server message: %s (%d)", t, t)
	}

	b := make([]byte, l-1)
	if _, err := io.ReadFull(c.transport, b); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(b, m); err != nil {
		return nil, err
	}

	// unwrap notice frames, return variable and state changes, skip over warnings
	if t == mysqlx.ServerMessages_NOTICE {
		f := m.(*mysqlx_notice.Frame)
		switch f.GetType() {
		case 1:
			m = new(mysqlx_notice.Warning)
			if err := proto.Unmarshal(f.Payload, m); err != nil {
				return nil, err
			}

			// TODO expose warnings?
			c.tracef("<== %T %v: %T %v", f, f, m, m)
			return c.readMessage(ctx)
		case 2:
			m = new(mysqlx_notice.SessionVariableChanged)
			if err := proto.Unmarshal(f.Payload, m); err != nil {
				return nil, err
			}
		case 3:
			m = new(mysqlx_notice.SessionStateChanged)
			if err := proto.Unmarshal(f.Payload, m); err != nil {
				return nil, err
			}
		default:
			return nil, bugf("conn.readMessage: unexpected notice frame type: %v", f)
		}

		if f.GetScope() != mysqlx_notice.Frame_LOCAL {
			return nil, bugf("conn.readMessage: unexpected notice frame scope: %v", f)
		}
	}

	c.tracef("<<< %T %v", m, m)
	return m, nil
}

// check interfaces
var (
	_ driver.Conn               = (*conn)(nil)
	_ driver.ConnBeginTx        = (*conn)(nil)
	_ driver.ConnPrepareContext = (*conn)(nil)
	_ driver.Execer             = (*conn)(nil)
	_ driver.ExecerContext      = (*conn)(nil)
	_ driver.Queryer            = (*conn)(nil)
	_ driver.QueryerContext     = (*conn)(nil)
	_ driver.Pinger             = (*conn)(nil)

	// TODO
	// _ driver.NamedValueChecker  = (*conn)(nil)
)
