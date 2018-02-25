// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
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

type connectionsMapType map[*conn]struct{}

func (cmt connectionsMapType) get() *conn {
	connectionsM.Lock()
	defer connectionsM.Unlock()

	if len(connectionsMap) != 1 {
		panic(fmt.Errorf("expected 1 connection, got %d", len(connectionsMap)))
	}
	for c := range connectionsMap {
		return c
	}
	panic("not reached")
}

var (
	connectionsProfile = pprof.NewProfile("github.com/AlekSi/mysqlx.connections")

	trackConnections bool
	connectionsM     sync.Mutex
	connectionsMap   = make(connectionsMapType)
)

func connectionOpened(c *conn) {
	connectionsProfile.Add(c, 1)

	if trackConnections {
		connectionsM.Lock()
		connectionsMap[c] = struct{}{}
		connectionsM.Unlock()
	}
}

func connectionClosed(c *conn) {
	connectionsProfile.Remove(c)

	if trackConnections {
		connectionsM.Lock()
		delete(connectionsMap, c)
		connectionsM.Unlock()
	}
}
