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

func TestAuthData(t *testing.T) {
	t.Parallel()

	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}

	assert.Equal(t, "\x00root\x00",
		string(authData("", "root", "", decode("434169533f3569721167252e59117a645a681500"))))
	assert.Equal(t, "world_x\x00root\x00",
		string(authData("world_x", "root", "", decode("434169533f3569721167252e59117a645a681500"))))
	assert.Equal(t, "\x00my_user\x00*C8B66ADB21E1E674249869852AEBC573DB7A5639",
		string(authData("", "my_user", "my_password", decode("434169533f3569721167252e59117a645a681500"))))
	assert.Equal(t, "world_x\x00my_user\x00*C8B66ADB21E1E674249869852AEBC573DB7A5639",
		string(authData("world_x", "my_user", "my_password", decode("434169533f3569721167252e59117a645a681500"))))
}
