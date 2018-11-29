// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeverityStringer(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "ERROR", SeverityError.String())
	assert.Equal(t, "FATAL", SeverityFatal.String())
	assert.Equal(t, "Severity 42", Severity(42).String())
}
