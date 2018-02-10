// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthDataSHA256(t *testing.T) {
	t.Parallel()

	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}

	assert.Equal(t, "\x00root\x0093AB01DB16BDA656A4164855EDFCBE4B5355A6AA726ED094EBBD517E836C31C2",
		string(authDataSHA256("", "root", "", decode("6c093446472736302a0b600f446e05212a6e7400"))))

	assert.Equal(t, "world_x\x00root\x0093AB01DB16BDA656A4164855EDFCBE4B5355A6AA726ED094EBBD517E836C31C2",
		string(authDataSHA256("world_x", "root", "", decode("6c093446472736302a0b600f446e05212a6e7400"))))

	assert.Equal(t, "\x00my_user\x007773D5FEE8462642B6141EACC3ED4F47104C1DDCDFBB4A0B5602452D4C73200D",
		string(authDataSHA256("", "my_user", "my_password", decode("6c093446472736302a0b600f446e05212a6e7400"))))

	assert.Equal(t, "world_x\000my_user\0007773D5FEE8462642B6141EACC3ED4F47104C1DDCDFBB4A0B5602452D4C73200D",
		string(authDataSHA256("world_x", "my_user", "my_password", decode("6c093446472736302a0b600f446e05212a6e7400"))))
}
