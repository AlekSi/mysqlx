// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"crypto/sha256"
	"fmt"
)

func scrambleSHA256(password string, authData []byte) []byte {
	// SHA256(password) XOR SHA256(SHA256(SHA256(password)) + authData)
	hash1 := sha256.Sum256([]byte(password))
	hash1h := sha256.Sum256(hash1[:])
	h := sha256.New()
	h.Write(hash1h[:])
	h.Write(authData)
	hash2 := h.Sum(nil)
	res := make([]byte, sha256.Size)
	for i := 0; i < sha256.Size; i++ {
		res[i] = hash1[i] ^ hash2[i]
	}
	return res
}

// https://dev.mysql.com/worklog/task/?id=10992
func authDataSHA256(database, username, password string, authData []byte) []byte {
	if len(authData) != 20 {
		return []byte(bugf("authDataSHA256: expected authData to has 20 bytes, got %d", len(authData)).Error())
	}

	res := database + "\x00" + username + "\x00"

	res += fmt.Sprintf("%X", scrambleSHA256(password, authData))
	return []byte(res)
}
