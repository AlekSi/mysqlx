package mysqlx

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/AlekSi/mysqlx/internal/mysqlx"
	"github.com/AlekSi/mysqlx/internal/mysqlx_connection"
	"github.com/AlekSi/mysqlx/internal/mysqlx_notice"
	"github.com/AlekSi/mysqlx/internal/mysqlx_resultset"
	"github.com/AlekSi/mysqlx/internal/mysqlx_session"
	"github.com/AlekSi/mysqlx/internal/mysqlx_sql"
)

var debugf = func(format string, v ...interface{}) {}

type conn struct {
	transport net.Conn

	closeOnce *sync.Once
	closeErr  error
}

func newConn(transport net.Conn) *conn {
	return &conn{
		transport: transport,

		closeOnce: new(sync.Once),
	}
}

func open(dataSource string) (*conn, error) {
	u, err := url.Parse(dataSource)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial(u.Scheme, u.Host)
	if err != nil {
		return nil, err
	}
	c := newConn(conn)

	database := strings.TrimPrefix(u.Path, "/")
	var username, password string
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	if err = c.negotiate(); err != nil {
		return nil, err
	}
	if err = c.auth(database, username, password); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *conn) negotiate() error {
	if err := writeMessage(c.transport, &mysqlx_connection.CapabilitiesGet{}); err != nil {
		return c.close(err)
	}
	m, err := readMessage(c.transport)
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

func (c *conn) auth(database, username, password string) error {
	// TODO use password

	mechName := "MYSQL41"
	if err := writeMessage(c.transport, &mysqlx_session.AuthenticateStart{
		MechName: &mechName,
	}); err != nil {
		return c.close(err)
	}

	m, err := readMessage(c.transport)
	if err != nil {
		return c.close(err)
	}
	cont := m.(*mysqlx_session.AuthenticateContinue)
	if len(cont.AuthData) != 20 {
		panic(len(cont.AuthData))
	}

	if err = writeMessage(c.transport, &mysqlx_session.AuthenticateContinue{
		AuthData: []byte(database + "\x00" + username + "\x00"),
	}); err != nil {
		return c.close(err)
	}

	m, err = readMessage(c.transport)
	if err != nil {
		return c.close(err)
	}
	_ = m.(*mysqlx_notice.SessionStateChanged)

	m, err = readMessage(c.transport)
	if err != nil {
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
	return c.closeErr
}

func (c *conn) Close() error {
	if err := writeMessage(c.transport, &mysqlx_connection.Close{}); err != nil {
		return c.close(err)
	}

	// read one next message, but do not check it is mysqlx.Ok
	if _, err := readMessage(c.transport); err != nil {
		return c.close(err)
	}

	return c.close(nil)
}

func (c *conn) Begin() (driver.Tx, error) {
	bugf("Begin not implemented yet")
	panic("not reached")
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	bugf("Prepare not implemented yet")
	panic("not reached")
}

func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	stmt := &mysqlx_sql.StmtExecute{
		Stmt: []byte(query),
	}
	for _, arg := range args {
		stmt.Args = append(stmt.Args, marshalValue(arg))
	}

	if err := writeMessage(c.transport, stmt); err != nil {
		return nil, c.close(err)
	}

	rows := rows{
		c:       c,
		columns: make([]mysqlx_resultset.ColumnMetaData, 0, 1),
		rows:    make(chan *mysqlx_resultset.Row, 1),
	}
	for {
		m, err := readMessage(c.transport)
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
		case *mysqlx_resultset.Row:
			rows.rows <- m
			go rows.runReader()
			return &rows, nil
		default:
			bugf("unhandled type %T", m)
		}
	}
}

func writeMessage(w io.Writer, m proto.Message) error {
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
		bugf("unhandled client message: %T %#v", m, m)
	}

	debugf(">>> %T %v", m, m)

	var head [5]byte
	binary.LittleEndian.PutUint32(head[:], uint32(len(b))+1)
	head[4] = byte(t)
	_, err = (&net.Buffers{head[:], b}).WriteTo(w)
	return err
}

func readMessage(r io.Reader) (proto.Message, error) {
	var head [5]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
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
		bugf("unhandled type of server message: %s (%d)", t, t)
	}

	b := make([]byte, l-1)
	if _, err := io.ReadFull(r, b); err != nil {
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
			debugf("<-- %T %v: %T %v", f, f, m, m)
			return readMessage(r)
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
			bugf("unexpected notice frame type: %v", f)
		}

		if f.GetScope() != mysqlx_notice.Frame_LOCAL {
			bugf("unexpected notice frame scope: %v", f)
		}
	}

	debugf("<<< %T %v", m, m)
	return m, nil
}

// check interfaces
var (
	_ driver.Conn    = (*conn)(nil)
	_ driver.Queryer = (*conn)(nil)

	// TODO
	// _ driver.ConnBeginTx        = (*conn)(nil)
	// _ driver.ConnPrepareContext = (*conn)(nil)
	// _ driver.Execer             = (*conn)(nil)
	// _ driver.ExecerContext      = (*conn)(nil)
	// _ driver.NamedValueChecker  = (*conn)(nil)
	// _ driver.Pinger             = (*conn)(nil)
	// _ driver.QueryerContext     = (*conn)(nil)
)
