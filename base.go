// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

const bug = true

func bugf(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...) + "\nPlease report this bug: https://github.com/AlekSi/mysqlx/issues\n"
	if bug {
		panic(msg)
	}
	return errors.New(msg)
}

// Severity represents Error severity level.
type Severity byte

const (
	// SeverityError indicates the current message sequence is aborted for the given error
	// and the session is ready for more.
	SeverityError Severity = 0

	// SeverityFatal indicates the client should not expect the server to continue handling any further messages
	// and should close the connection.
	SeverityFatal Severity = 1
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "ERROR"
	case SeverityFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("Severity %d", s)
	}
}

// Error represents MySQL X Protocol error message.
// It's not used for transport-level errors.
type Error struct {
	Severity Severity
	Code     uint32
	SQLState string
	Msg      string
}

func (e *Error) Error() string {
	// format of mysql and mysqlsh client programs
	return fmt.Sprintf("%s %d (%s): %s", e.Severity, e.Code, e.SQLState, e.Msg)
}

type execResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (r execResult) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}

func (r execResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

// check interfaces
var (
	_ fmt.Stringer  = SeverityError
	_ error         = (*Error)(nil)
	_ driver.Result = execResult{}
)
