// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/AlekSi/mysqlx/internal/proto/mysqlx"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_connection"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_datatypes"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_notice"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_resultset"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_session"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_sql"
)

// TODO make this configurable?
// It should not be less then 1.
const rowsCap = 1

// https://github.com/mysql/mysql-server/blob/mysql-5.7.19/rapid/plugin/x/mysqlxtest_src/password_hasher.cc
// https://github.com/mysql/mysql-server/blob/mysql-5.7.19/rapid/plugin/x/mysqlxtest_src/mysql41_hash.cc
func scramble(password string, authData []byte) []byte {
	hash1 := sha1.Sum([]byte(password))
	hash2 := sha1.Sum([]byte(hash1[:]))

	h := sha1.New()
	h.Write(authData)
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

	res += fmt.Sprintf("*%X", scramble(password, authData))
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
	local, remote := transport.LocalAddr().String(), transport.RemoteAddr().String()
	traceF("+++ connection created: %s->%s", local, remote)
	connectionCreated(local, remote)

	return &conn{
		transport: transport,
		tracef:    traceF,
	}
}

func open(ctx context.Context, connector *Connector) (*conn, error) {
	conn, err := new(net.Dialer).DialContext(ctx, "tcp", connector.hostPort())
	if err != nil {
		return nil, err
	}
	c := newConn(conn, connector.Trace)
	defer func() {
		if err != nil {
			c.close(err)
		}
	}()

	if err = c.negotiate(ctx); err != nil {
		return nil, err
	}
	if err = c.authenticate(ctx, connector.AuthMethod, connector.Database, connector.Username, connector.Password); err != nil {
		return nil, err
	}

	// set session variables
	keys := make([]string, 0, len(connector.SessionVariables))
	for k := range connector.SessionVariables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		nv := []driver.NamedValue{{
			Ordinal: 1,
			Value:   connector.SessionVariables[k],
		}}
		if _, err = c.ExecContext(ctx, "SET SESSION "+k+" = ?", nv); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *conn) negotiate(ctx context.Context) error {
	if err := c.writeMessage(ctx, &mysqlx_connection.CapabilitiesGet{}); err != nil {
		return err
	}
	m, err := c.readMessage(ctx)
	if err != nil {
		return err
	}

	var tlsFound bool
	for _, cap := range m.(*mysqlx_connection.Capabilities).Capabilities {
		if cap.GetName() == "tls" {
			tlsFound = true
		}
	}

	// enable TLS if possible
	if tlsFound {
		cap := &mysqlx_connection.CapabilitiesSet{
			Capabilities: &mysqlx_connection.Capabilities{
				Capabilities: []*mysqlx_connection.Capability{{
					Name: proto.String("tls"),
					Value: &mysqlx_datatypes.Any{
						Type: mysqlx_datatypes.Any_SCALAR.Enum(),
						Scalar: &mysqlx_datatypes.Scalar{
							Type:  mysqlx_datatypes.Scalar_V_BOOL.Enum(),
							VBool: proto.Bool(true),
						},
					},
				}},
			},
		}
		if err := c.writeMessage(ctx, cap); err != nil {
			return err
		}
		if _, err := c.readMessage(ctx); err != nil {
			return err
		}
		// FIXME
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		tlsConn := tls.Client(c.transport, tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			tlsConn.Close()
			return err
		}
		c.transport = tlsConn
	}

	return nil
}

func (c *conn) authenticate(ctx context.Context, method AuthMethod, database, username, password string) error {
	/*
		if err := c.writeMessage(ctx, &mysqlx_connection.CapabilitiesGet{}); err != nil {
			return err
		}
		m, err := c.readMessage(ctx)
		if err != nil {
			return err
		}

		var tlsEnabled, sha256Found, mysql41Found, plainFound bool
			for _, cap := range m.(*mysqlx_connection.Capabilities).Capabilities {
				switch cap.GetName() {
				case "tls":
					tlsEnabled = cap.GetValue().GetScalar().GetVBool()
				case "authentication.mechanisms":
					for _, value := range cap.Value.Array.Value {
						s := string(value.Scalar.VString.Value)
						switch s {
						case "SHA256_MEMORY":
							// FIXME
							// sha256Found = true
						case "MYSQL41":
							mysql41Found = true
						case "PLAIN":
							plainFound = true
						}
					}
				}
			}

			err = fmt.Errorf("can't authenticate")
			if err != nil && sha256Found {
				err = c.authSHA256(ctx, database, username, password)
			}
			if err != nil && mysql41Found {
				err = c.authMySQL41(ctx, database, username, password)
			}
			if err != nil && plainFound && tlsEnabled {
				err = c.authPlain(ctx, database, username, password)
			}
			return err
	*/

	switch method {
	case AuthPlain:
		return c.authPlain(ctx, database, username, password)
	case AuthMySQL41:
		return c.authMySQL41(ctx, database, username, password)
	default:
		return fmt.Errorf("unexpected authentication method %q", method)
	}
}

func (c *conn) authSHA256(ctx context.Context, database, username, password string) error {
	if err := c.writeMessage(ctx, &mysqlx_session.AuthenticateStart{
		MechName: proto.String("SHA256_MEMORY"),
	}); err != nil {
		return err
	}

	m, err := c.readMessage(ctx)
	if err != nil {
		return err
	}
	cont := m.(*mysqlx_session.AuthenticateContinue)

	if err = c.writeMessage(ctx, &mysqlx_session.AuthenticateContinue{
		AuthData: authDataSHA256(database, username, password, cont.AuthData),
	}); err != nil {
		return err
	}

	if m, err = c.readMessage(ctx); err != nil {
		return err
	}
	switch m := m.(type) {
	case *mysqlx.Error:
		severity := Severity(m.GetSeverity())
		return &Error{
			Severity: severity,
			Code:     m.GetCode(),
			SQLState: m.GetSqlState(),
			Msg:      m.GetMsg(),
		}

	case *mysqlx_notice.SessionStateChanged:
	default:
		bugf("conn.authSHA256: unhandled type %T", m)
	}

	if m, err = c.readMessage(ctx); err != nil {
		return err
	}
	_ = m.(*mysqlx_session.AuthenticateOk)

	return nil
}

func (c *conn) authMySQL41(ctx context.Context, database, username, password string) error {
	if err := c.writeMessage(ctx, &mysqlx_session.AuthenticateStart{
		MechName: proto.String("MYSQL41"),
	}); err != nil {
		return err
	}

	m, err := c.readMessage(ctx)
	if err != nil {
		return err
	}
	cont := m.(*mysqlx_session.AuthenticateContinue)

	if err = c.writeMessage(ctx, &mysqlx_session.AuthenticateContinue{
		AuthData: authData(database, username, password, cont.AuthData),
	}); err != nil {
		return err
	}

	if m, err = c.readMessage(ctx); err != nil {
		return err
	}
	switch m := m.(type) {
	case *mysqlx.Error:
		severity := Severity(m.GetSeverity())
		return &Error{
			Severity: severity,
			Code:     m.GetCode(),
			SQLState: m.GetSqlState(),
			Msg:      m.GetMsg(),
		}

	case *mysqlx_notice.SessionStateChanged:
	default:
		bugf("conn.authMySQL41: unhandled type %T", m)
	}

	if m, err = c.readMessage(ctx); err != nil {
		return err
	}
	_ = m.(*mysqlx_session.AuthenticateOk)

	return nil
}

func (c *conn) authPlain(ctx context.Context, database, username, password string) error {
	if err := c.writeMessage(ctx, &mysqlx_session.AuthenticateStart{
		MechName: proto.String("PLAIN"),
		AuthData: []byte(database + "\x00" + username + "\x00" + password),
	}); err != nil {
		return err
	}

	var m proto.Message
	var err error
	if m, err = c.readMessage(ctx); err != nil {
		return err
	}
	switch m := m.(type) {
	case *mysqlx.Error:
		severity := Severity(m.GetSeverity())
		return &Error{
			Severity: severity,
			Code:     m.GetCode(),
			SQLState: m.GetSqlState(),
			Msg:      m.GetMsg(),
		}

	case *mysqlx_notice.SessionStateChanged:
	default:
		bugf("conn.authPlain: unhandled type %T", m)
	}

	if m, err = c.readMessage(ctx); err != nil {
		return err
	}
	switch m := m.(type) {
	case *mysqlx_session.AuthenticateOk:
	default:
		bugf("conn.authPlain: unhandled type %T", m)
	}

	return nil
}

func (c *conn) close(err error) error {
	c.closeOnce.Do(func() {
		local, remote := c.transport.LocalAddr().String(), c.transport.RemoteAddr().String()
		c.closeErr = err
		e := c.transport.Close()
		if c.closeErr == nil {
			c.closeErr = e
		}
		connectionClosed(local, remote)
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
	var chars []string
	switch sql.IsolationLevel(opts.Isolation) {
	case sql.LevelDefault:
		// nothing
	case sql.LevelReadUncommitted:
		chars = append(chars, "ISOLATION LEVEL READ UNCOMMITTED")
	case sql.LevelReadCommitted:
		chars = append(chars, "ISOLATION LEVEL READ COMMITTED")
	case sql.LevelRepeatableRead:
		chars = append(chars, "ISOLATION LEVEL REPEATABLE READ")
	case sql.LevelSnapshot:
		// special handling below
	case sql.LevelSerializable:
		chars = append(chars, "ISOLATION LEVEL SERIALIZABLE")
	default:
		return nil, bugf("conn.BeginTx: isolation level %d is not supported yet", opts.Isolation)
	}
	if opts.ReadOnly {
		chars = append(chars, "READ ONLY")
	}
	if chars != nil {
		q := "SET TRANSACTION " + strings.Join(chars, ", ")
		if _, err := c.ExecContext(ctx, q, nil); err != nil {
			return nil, err
		}
	}

	q := "START TRANSACTION"
	if sql.IsolationLevel(opts.Isolation) == sql.LevelSnapshot {
		q += " WITH CONSISTENT SNAPSHOT"
	}
	if _, err := c.ExecContext(ctx, q, nil); err != nil {
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

		case *mysqlx_notice.Warning:
			// TODO expose warnings?
			continue

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

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
// If Ping returns driver.ErrBadConn, sql.DB.Ping and sql.DB.PingContext will remove the Conn from pool.
func (c *conn) Ping(ctx context.Context) error {
	if _, err := c.ExecContext(ctx, "SELECT 'ping'", nil); err != nil {
		return driver.ErrBadConn
	}
	return nil
}

// CheckNamedValue is called before passing arguments to the driver
// and is called in place of any ColumnConverter. CheckNamedValue must do type
// validation and conversion as appropriate for the driver.
func (c *conn) CheckNamedValue(arg *driver.NamedValue) error {
	if arg.Name != "" {
		return bugf("conn.CheckNamedValue: %q - named values are not supported yet", arg.Name)
	}

	// pass everything to datatypes handling (marshalValue and unmarshalValue)
	return nil
}

// writeMessage writes one protocol message, returns low-level error if any.
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
	case *mysqlx_connection.CapabilitiesSet:
		t = mysqlx.ClientMessages_CON_CAPABILITIES_SET
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
	_, err = (&net.Buffers{head[:], b}).WriteTo(c.transport) // use writev(2) if available
	return err
}

// ReadMessage reads and returns one next protocol message, or low-level error.
// Notices are unwrapped: SessionVariableChanged, SessionStateChanged, and Warning are returned,
// and raw Frame is never returned.
// TODO un-export (currently required for mitm-proxy)
func ReadMessage(r io.Reader) (proto.Message, []byte, error) {
	var head [5]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return nil, nil, err
	}

	buf := make([]byte, binary.LittleEndian.Uint32(head[:])+4)
	copy(buf, head[:])
	if _, err := io.ReadFull(r, buf[5:]); err != nil {
		return nil, nil, err
	}

	t := mysqlx.ServerMessages_Type(buf[4])
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
		return nil, nil, bugf("conn.readMessage: unhandled type of server message: %s (%d)", t, t)
	}

	if err := proto.Unmarshal(buf[5:], m); err != nil {
		return nil, buf, fmt.Errorf("conn.readMessage: %s", err)
	}

	// unwrap notice frames, return variable and state changes, skip over warnings
	if t == mysqlx.ServerMessages_NOTICE {
		f := m.(*mysqlx_notice.Frame)
		switch mysqlx_notice.Frame_Type(f.GetType()) {
		case mysqlx_notice.Frame_WARNING:
			m = new(mysqlx_notice.Warning)
			if err := proto.Unmarshal(f.Payload, m); err != nil {
				return nil, nil, err
			}
		case mysqlx_notice.Frame_SESSION_VARIABLE_CHANGED:
			m = new(mysqlx_notice.SessionVariableChanged)
			if err := proto.Unmarshal(f.Payload, m); err != nil {
				return nil, nil, err
			}
		case mysqlx_notice.Frame_SESSION_STATE_CHANGED:
			m = new(mysqlx_notice.SessionStateChanged)
			if err := proto.Unmarshal(f.Payload, m); err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, bugf("conn.readMessage: unexpected notice frame type: %v", f)
		}

		if f.GetScope() != mysqlx_notice.Frame_LOCAL {
			return nil, nil, bugf("conn.readMessage: unexpected notice frame scope: %v", f)
		}
	}

	return m, buf, nil
}

func (c *conn) readMessage(ctx context.Context) (proto.Message, error) {
	deadline, _ := ctx.Deadline()
	if err := c.transport.SetReadDeadline(deadline); err != nil {
		return nil, err
	}
	m, _, err := ReadMessage(c.transport)
	if err != nil {
		return nil, err
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
	_ driver.NamedValueChecker  = (*conn)(nil)
	// _ driver.Connector          = (*conn)(nil)
)
