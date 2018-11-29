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

type decimalInput struct {
	b        []byte
	expected string
	valid    bool
}

var decimalInputs = []decimalInput{
	// valid values from TestQueryData
	{[]byte{0x00, 0x9c}, "9", true},
	{[]byte{0x00, 0x9d}, "-9", true},
	{[]byte{0x00, 0x12, 0xc0}, "12", true},
	{[]byte{0x00, 0x12, 0xd0}, "-12", true},
	{[]byte{0x01, 0x09, 0xc0}, "0.9", true},
	{[]byte{0x01, 0x09, 0xd0}, "-0.9", true},
	{[]byte{0x03, 0x12, 0x34, 0x0c}, "12.340", true},
	{[]byte{0x03, 0x12, 0x34, 0x0d}, "-12.340", true},
	{[]byte{0x04, 0x12, 0x34, 0x01, 0xc0}, "12.3401", true},
	{[]byte{0x04, 0x12, 0x34, 0x01, 0xd0}, "-12.3401", true},

	{nil, "", false},
	{[]byte{}, "", false},
	{[]byte{0x00}, "", false},
	{[]byte{0x00, 0x00}, "", false},
	{[]byte{0x30, 0x30}, "", false},
	{[]byte{0xff, 0xff}, "", false},
	{[]byte{0x30, 0x0a, 0x30}, "", false},
}

func TestUnmarshalDecimal(t *testing.T) {
	t.Parallel()

	for _, input := range decimalInputs {
		d, err := unmarshalDecimal(input.b)
		assert.Equal(t, input.valid, err == nil, "%s", err)
		assert.Equal(t, input.expected, d)
	}
}

var sink interface{}

func BenchmarkUnmarshalDecimal(b *testing.B) {
	for _, input := range decimalInputs {
		b.Run(input.expected, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sink, sink = unmarshalDecimal(input.b)
			}
		})
	}
}
