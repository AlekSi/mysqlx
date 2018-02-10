// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gofuzz

package mysqlx

func FuzzUnmarshalDecimal(data []byte) int {
	_, err := unmarshalDecimal(data)
	if err != nil {
		return 0
	}
	return 1
}
