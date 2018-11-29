// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"fmt"
	"runtime/pprof"
	"sync"
)

type trackedConnections struct {
	l sync.Mutex
	m map[*conn]struct{}
}

func (tc *trackedConnections) get() *conn {
	tc.l.Lock()
	defer tc.l.Unlock()

	if len(tc.m) != 1 {
		panic(fmt.Errorf("expected 1 connection, got %d", len(tc.m)))
	}
	for c := range tc.m {
		return c
	}
	panic("not reached")
}

var (
	connectionsProfile = pprof.NewProfile("github.com/AlekSi/mysqlx.connections")

	// initialized in tracking_test.go for testing only
	testConnections *trackedConnections
)

func connectionOpened(c *conn) {
	connectionsProfile.Add(c, 1)

	if testConnections != nil {
		testConnections.l.Lock()
		testConnections.m[c] = struct{}{}
		testConnections.l.Unlock()
	}
}

func connectionClosed(c *conn) {
	connectionsProfile.Remove(c)

	if testConnections != nil {
		testConnections.l.Lock()
		delete(testConnections.m, c)
		testConnections.l.Unlock()
	}
}
