// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDataSource(t *testing.T) {
	t.Parallel()

	for s, expected := range map[string]*DataSource{
		"mysqlx://my_user:my_password@127.0.0.1:33060/world_x?time_zone=UTC": {
			Host:             "127.0.0.1",
			Port:             33060,
			Database:         "world_x",
			Username:         "my_user",
			Password:         "my_password",
			SessionVariables: map[string]string{"time_zone": "UTC"},
		},
	} {
		t.Run(s, func(t *testing.T) {
			t.Parallel()

			u, err := url.Parse(s)
			require.NoError(t, err)
			actual, err := ParseDataSource(u)
			require.NoError(t, err)
			actual.Trace = nil // "Func values are deeply equal if both are nil; otherwise they are not deeply equal."
			assert.Equal(t, expected, actual)
			assert.Equal(t, u, actual.URL())
		})
	}
}