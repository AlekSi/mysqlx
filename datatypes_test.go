package mysqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type decimalPair struct {
	b        []byte
	expected string
}

// values from TestQueryData
var decimalPairs = []decimalPair{
	{[]byte{0x04, 0x12, 0x34, 0x01, 0xc0}, "12.3401"},
	{[]byte{0x04, 0x12, 0x34, 0x01, 0xd0}, "-12.3401"},
	{[]byte{0x03, 0x12, 0x34, 0x0c}, "12.340"},
	{[]byte{0x03, 0x12, 0x34, 0x0d}, "-12.340"},
	{[]byte{0x00, 0x12, 0xc0}, "12"},
	{[]byte{0x00, 0x12, 0xd0}, "-12"},
	{[]byte{0x00, 0x9c}, "9"},
	{[]byte{0x0, 0x9d}, "-9"},
	{[]byte{0x1, 0x9, 0xc0}, "0.9"},
	{[]byte{0x1, 0x9, 0xd0}, "-0.9"},
}

func TestDecimal(t *testing.T) {
	for _, p := range decimalPairs {
		assert.Equal(t, p.expected, unmarshalDecimal(p.b))
	}
}

var sink string

func BenchmarkDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, p := range decimalPairs {
			sink = unmarshalDecimal(p.b)
		}
	}
}
