// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

// https://github.com/mysql/mysql-server/blob/mysql-5.7.19/rapid/plugin/x/mysqlxtest_src/password_hasher.cc
// https://github.com/mysql/mysql-server/blob/mysql-5.7.19/rapid/plugin/x/mysqlxtest_src/mysql41_hash.cc
func scrambleMySQL41(password string, authData []byte) []byte {
	hash1 := sha1.Sum([]byte(password))
	hash2 := sha1.Sum(hash1[:])

	h := sha1.New()
	h.Write(authData)
	h.Write(hash2[:])
	res := h.Sum(nil)

	for i := range res {
		res[i] ^= hash1[i]
	}
	return res[:]
}

func authDataMySQL41(database, username, password string, authData []byte) ([]byte, error) {
	if len(authData) != 20 {
		return nil, fmt.Errorf("authDataMySQL41: expected authData to has 20 bytes, got %d", len(authData))
	}

	res := database + "\x00" + username + "\x00"
	if password == "" {
		return []byte(res), nil
	}

	res += fmt.Sprintf("*%X", scrambleMySQL41(password, authData))
	return []byte(res), nil
}

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
func authDataSHA256(database, username, password string, authData []byte) ([]byte, error) {
	if len(authData) != 20 {
		return nil, fmt.Errorf("authDataSHA256: expected authData to has 20 bytes, got %d", len(authData))
	}

	res := database + "\x00" + username + "\x00"

	res += fmt.Sprintf("%X", scrambleSHA256(password, authData))
	return []byte(res), nil
}
