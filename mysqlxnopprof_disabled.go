// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build !mysqlxnopprof

package mysqlx

import (
	"fmt"
	"runtime/pprof"
)

var (
	connectionsProfile = pprof.NewProfile("github.com/AlekSi/mysqlx.connections")
)

func connectionCreated(local, remote string) {
	key := fmt.Sprintf("%s-%s", local, remote)
	connectionsProfile.Add(key, 1)
}

func connectionClosed(local, remote string) {
	key := fmt.Sprintf("%s-%s", local, remote)
	connectionsProfile.Remove(key)
}
