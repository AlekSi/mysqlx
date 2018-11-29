// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mysqlx

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_datatypes"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_resultset"
)

var btoa = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func unmarshalDecimal(value []byte) (string, error) {
	addToGoFuzzCorpus("unmarshalDecimal", value)

	if len(value) < 2 {
		return "", fmt.Errorf("unmarshalDecimal: failed to parse decimal %#v", value)
	}

	sign := value[len(value)-1]
	var s string
	for _, b := range value[1 : len(value)-1] {
		h := (b >> 4) & 0x0f
		l := b & 0x0f
		if h > 9 || l > 9 {
			return "", fmt.Errorf("unmarshalDecimal: failed to parse decimal %#v", value)
		}
		s += btoa[h] + btoa[l]
	}
	if sign != 0xd0 && sign != 0xc0 {
		h := (sign >> 4) & 0x0f
		if h > 9 {
			return "", fmt.Errorf("unmarshalDecimal: failed to parse decimal %#v", value)
		}
		s += btoa[h]
		sign <<= 4
	}
	if scale := int(value[0]); scale != 0 {
		if scale >= len(s) {
			return "", fmt.Errorf("unmarshalDecimal: failed to parse decimal %#v", value)
		}
		s = s[:len(s)-scale] + "." + s[len(s)-scale:]
	}
	switch sign {
	case 0xd0:
		return "-" + s, nil
	case 0xc0:
		return s, nil
	default:
		return "", fmt.Errorf("unmarshalDecimal: failed to parse decimal %#v", value)
	}
}

func unmarshalValue(value []byte, column *mysqlx_resultset.ColumnMetaData) (driver.Value, error) {
	// NULL -> nil, ignore type
	if len(value) == 0 {
		return nil, nil
	}

	// in order of ColumnMetaData_FieldType values
	switch *column.Type {
	case mysqlx_resultset.ColumnMetaData_SINT:
		// TINY, SHORT, INT24, INT, LONGLONG
		i64, n := binary.Varint(value)
		if n != len(value) {
			return nil, bugf("unmarshalValue: failed to decode %#v as SINT", value)
		}
		return i64, nil

	case mysqlx_resultset.ColumnMetaData_UINT:
		// TINY UNSIGNED, SHORT UNSIGNED, INT24 UNSIGNED, INT UNSIGNED, LONGLONG UNSIGNED, YEAR
		u64, n := binary.Uvarint(value)
		if n != len(value) {
			return nil, bugf("unmarshalValue: failed to decode %#v as UINT", value)
		}
		return u64, nil

	case mysqlx_resultset.ColumnMetaData_DOUBLE:
		// DOUBLE
		u64, err := proto.NewBuffer(value).DecodeFixed64()
		if err != nil {
			return nil, bugf("unmarshalValue: failed to decode %#v as DOUBLE: %s", value, err)
		}
		return math.Float64frombits(u64), nil

	case mysqlx_resultset.ColumnMetaData_FLOAT:
		// FLOAT
		u64, err := proto.NewBuffer(value).DecodeFixed32()
		if err != nil {
			return nil, bugf("unmarshalValue: failed to decode %#v as FLOAT: %s", value, err)
		}
		return float64(math.Float32frombits(uint32(u64))), nil

	case mysqlx_resultset.ColumnMetaData_BYTES:
		// VARCHAR, CHAR, GEOMETRY (and also NULL, but we handle it separately)
		return string(value[:len(value)-1]), nil // trim last 0x00

	case mysqlx_resultset.ColumnMetaData_TIME:
		// TIME
		// FIXME convert to time.Duration? what about range?
		// and time.Duration is not a driver.Value!
		return nil, bugf("unmarshalValue: unhandled TIME %s, value %#v", column.Type, value)

	case mysqlx_resultset.ColumnMetaData_DATETIME:
		// DATE, DATETIME, TIMESTAMP
		// year, month and day are mandatory, other parts are optional
		r := bytes.NewReader(value)
		year, _ := binary.ReadUvarint(r)
		month, _ := binary.ReadUvarint(r)
		day, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, bugf("unmarshalValue: failed to decode %#v as DATETIME: %s", value, err)
		}
		hour, _ := binary.ReadUvarint(r)
		min, _ := binary.ReadUvarint(r)
		sec, _ := binary.ReadUvarint(r)
		usec, _ := binary.ReadUvarint(r)
		return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), int(usec)*1000, time.UTC), nil

	case mysqlx_resultset.ColumnMetaData_SET:
		// SET
		return nil, bugf("unmarshalValue: unhandled SET %s, value %#v", column.Type, value)

	case mysqlx_resultset.ColumnMetaData_ENUM:
		// ENUM
		return string(value[:len(value)-1]), nil // trim last 0x00

	case mysqlx_resultset.ColumnMetaData_BIT:
		// BIT
		return nil, bugf("unmarshalValue: unhandled BIT %s, value %#v", column.Type, value)

	case mysqlx_resultset.ColumnMetaData_DECIMAL:
		// DECIMAL
		return unmarshalDecimal(value)

	default:
		return nil, bugf("unmarshalValue: unhandled type %s, value %#v", column.Type, value)
	}
}

func marshalValue(value driver.Value) (*mysqlx_datatypes.Any, error) {
	// Due to conn.CheckNamedValue passing on everything, value there can be of any type, not only of the one of
	// standard driver.Value types. We should handle everything ourselves.

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.IsValid() {
			value = v.Interface()
		} else {
			value = nil
		}
	}

	// nil -> NULL
	if value == nil {
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type: mysqlx_datatypes.Scalar_V_NULL.Enum(),
			},
		}, nil
	}

	// in order of Scalar_Type values
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		// SINT
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type:       mysqlx_datatypes.Scalar_V_SINT.Enum(),
				VSignedInt: proto.Int64(reflect.ValueOf(value).Int()),
			},
		}, nil

	case uint, uint8, uint16, uint32, uint64:
		// UINT
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type:         mysqlx_datatypes.Scalar_V_UINT.Enum(),
				VUnsignedInt: proto.Uint64(reflect.ValueOf(value).Uint()),
			},
		}, nil

	case float64:
		// DOUBLE
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type:    mysqlx_datatypes.Scalar_V_DOUBLE.Enum(),
				VDouble: &value,
			},
		}, nil

	case float32:
		// FLOAT
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type:   mysqlx_datatypes.Scalar_V_FLOAT.Enum(),
				VFloat: &value,
			},
		}, nil

	case bool:
		// BOOL
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type:  mysqlx_datatypes.Scalar_V_BOOL.Enum(),
				VBool: &value,
			},
		}, nil

	case string:
		// STRING
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type: mysqlx_datatypes.Scalar_V_STRING.Enum(),
				VString: &mysqlx_datatypes.Scalar_String{
					Value: []byte(value),
				},
			},
		}, nil

	case time.Time:
		s := value.Format("2006-01-02 15:04:05.999999999")
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type: mysqlx_datatypes.Scalar_V_OCTETS.Enum(),
				VOctets: &mysqlx_datatypes.Scalar_Octets{
					Value: []byte(s),
				},
			},
		}, nil

	case uintptr:
		return nil, fmt.Errorf("marshalValue: unhandled type %T, value %#v", value, value)

	default:
		return nil, bugf("marshalValue: unhandled type %T, value %#v", value, value)
	}
}
