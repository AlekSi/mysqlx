// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
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
)

type AuthMethod string

const (
	AuthPlain   AuthMethod = "PLAIN"
	AuthMySQL41 AuthMethod = "MYSQL41"
)

// Trace function signature.
type TraceFunc func(format string, v ...interface{})

// noTrace is a trace functions which does nothing.
// TODO check it is inlined and eliminated by compiler.
func noTrace(string, ...interface{}) {}

type DataSource struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string

	AuthMethod AuthMethod

	SessionVariables map[string]string

	Trace TraceFunc
}

func (ds *DataSource) Connect(ctx context.Context) (driver.Conn, error) {
	return open(ctx, ds)
}

func (ds *DataSource) Driver() driver.Driver {
	return Driver
}

func ParseDataSource(u *url.URL) (*DataSource, error) {
	if u.Scheme != "mysqlx" {
		return nil, fmt.Errorf("unexpected scheme %s", u.Scheme)
	}
	ds := &DataSource{
		Host:     u.Hostname(),
		Database: strings.TrimPrefix(u.Path, "/"),
		Trace:    noTrace,
	}

	// set port if given
	if p := u.Port(); p != "" {
		pp, err := strconv.ParseUint(p, 10, 16)
		if err != nil {
			return nil, err
		}
		ds.Port = uint16(pp)
	}

	// set username and password if they are given
	if u.User != nil {
		ds.Username = u.User.Username()
		ds.Password, _ = u.User.Password()
	}

	for k, vs := range u.Query() {
		if len(vs) != 1 {
			return nil, fmt.Errorf("%d values given for session variable %s: %v", len(vs), k, vs)
		}
		v := vs[0]

		// set session variables
		if !strings.HasPrefix(k, "_") {
			if ds.SessionVariables == nil {
				ds.SessionVariables = make(map[string]string)
			}
			ds.SessionVariables[k] = v
			continue
		}

		switch k {
		case "_auth-method":
			switch v {
			case string(AuthPlain):
				ds.AuthMethod = AuthPlain
			case string(AuthMySQL41):
				ds.AuthMethod = AuthMySQL41
			default:
				return nil, fmt.Errorf("unexpected value for %q: %q", k, v)
			}

		default:
			return nil, fmt.Errorf("unexpected parameter %q", k)
		}
	}

	return ds, nil
}

func (ds *DataSource) hostPort() string {
	return net.JoinHostPort(ds.Host, strconv.FormatUint(uint64(ds.Port), 10))
}

func (ds *DataSource) URL() *url.URL {
	u := &url.URL{
		Scheme: "mysqlx",
		Host:   ds.hostPort(),
		Path:   "/" + ds.Database,
	}

	if ds.Username != "" {
		u.User = url.UserPassword(ds.Username, ds.Password)
	}

	q := make(url.Values)
	if ds.AuthMethod != "" {
		q.Set("_auth-method", string(ds.AuthMethod))
	}

	for k, v := range ds.SessionVariables {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()
	return u
}

// check interfaces
var (
	_ driver.Connector = (*DataSource)(nil)
)
