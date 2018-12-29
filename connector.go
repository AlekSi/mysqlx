// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AuthMethod string

const (
	AuthPlain   AuthMethod = "PLAIN"
	AuthMySQL41 AuthMethod = "MYSQL41"
)

// noTrace is a trace functions which does nothing.
func noTrace(string, ...interface{}) {}

// Connector implements database/sql/driver.Connector interface.
type Connector struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string

	AuthMethod  AuthMethod
	DialTimeout time.Duration

	SessionVariables map[string]string

	Trace func(format string, v ...interface{}) // may be nil
}

// Connect returns a new connection to the database.
//
// The provided context.Context is for dialing purposes only
// (see net.DialContext) and is not used for other purposes.
//
// The returned connection must be used only by one goroutine at a time.
func (connector *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return open(ctx, connector)
}

// Driver returns the underlying Driver of the Connector,
// mainly to maintain compatibility with the Driver method on sql.DB.
func (connector *Connector) Driver() driver.Driver {
	return Driver
}

// ParseDataSource returns Connector for given data source.
func ParseDataSource(dataSource string) (*Connector, error) {
	u, err := url.Parse(dataSource)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "mysqlx" {
		return nil, fmt.Errorf("unexpected scheme %s", u.Scheme)
	}
	connector := &Connector{
		Host:     u.Hostname(),
		Database: strings.TrimPrefix(u.Path, "/"),
	}

	// set port if given
	if p := u.Port(); p != "" {
		pp, err := strconv.ParseUint(p, 10, 16)
		if err != nil {
			return nil, err
		}
		connector.Port = uint16(pp)
	}

	// set username and password if they are given
	if u.User != nil {
		connector.Username = u.User.Username()
		connector.Password, _ = u.User.Password()
	}

	for k, vs := range u.Query() {
		if len(vs) != 1 {
			return nil, fmt.Errorf("%d values given for session variable %s: %v", len(vs), k, vs)
		}
		v := vs[0]

		// set session variables
		if !strings.HasPrefix(k, "_") {
			if connector.SessionVariables == nil {
				connector.SessionVariables = make(map[string]string)
			}
			connector.SessionVariables[k] = v
			continue
		}

		switch k {
		case "_auth-method":
			switch v {
			case string(AuthPlain):
				connector.AuthMethod = AuthPlain
			case string(AuthMySQL41):
				connector.AuthMethod = AuthMySQL41
			default:
				return nil, fmt.Errorf("unexpected value for %q: %q", k, v)
			}

		case "_dial-timeout":
			connector.DialTimeout, err = time.ParseDuration(v)
			if err != nil {
				dt, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("unexpected value for %q: %q", k, v)
				}
				connector.DialTimeout = time.Duration(dt) * time.Second
			}

		default:
			return nil, fmt.Errorf("unexpected parameter %q", k)
		}
	}

	return connector, nil
}

func (connector *Connector) hostPort() string {
	return net.JoinHostPort(connector.Host, strconv.FormatUint(uint64(connector.Port), 10))
}

// URL returns data source as an URL.
func (connector *Connector) URL() *url.URL {
	u := &url.URL{
		Scheme: "mysqlx",
		Host:   connector.hostPort(),
		Path:   "/" + connector.Database,
	}

	if connector.Username != "" {
		u.User = url.UserPassword(connector.Username, connector.Password)
	}

	q := make(url.Values)
	if connector.AuthMethod != "" {
		q.Set("_auth-method", string(connector.AuthMethod))
	}

	for k, v := range connector.SessionVariables {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()
	return u
}

// check interfaces
var (
	_ driver.Connector = (*Connector)(nil)
)
