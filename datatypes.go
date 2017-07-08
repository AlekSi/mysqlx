package mysqlx

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/AlekSi/mysqlx/internal/mysqlx_datatypes"
	"github.com/AlekSi/mysqlx/internal/mysqlx_resultset"
)

func unmarshalValue(value []byte, column *mysqlx_resultset.ColumnMetaData) (driver.Value, error) {
	// NULL -> nil, ignore type
	if len(value) == 0 {
		return nil, nil
	}

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
		return int64(u64), nil

	case mysqlx_resultset.ColumnMetaData_BYTES:
		// VARCHAR, CHAR, GEOMETRY (and also NULL, but we handle it separately)
		return string(value[:len(value)-1]), nil // trim last 0x00

	// case mysqlx_resultset.ColumnMetaData_TIME:
	// TIME
	// FIXME convert to time.Duration? what about range?
	// and time.Duration is not a driver.Value!

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

	default:
		return nil, bugf("unmarshalValue: unhandled type %s", column.Type)
	}
}

func marshalValue(value driver.Value) (*mysqlx_datatypes.Any, error) {
	// nil -> NULL
	if value == nil {
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type: mysqlx_datatypes.Scalar_V_NULL.Enum(),
			},
		}, nil
	}

	switch value := value.(type) {
	case int64:
		return &mysqlx_datatypes.Any{
			Type: mysqlx_datatypes.Any_SCALAR.Enum(),
			Scalar: &mysqlx_datatypes.Scalar{
				Type:       mysqlx_datatypes.Scalar_V_SINT.Enum(),
				VSignedInt: proto.Int64(value),
			},
		}, nil

	case string:
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

	default:
		return nil, bugf("marshalValue: unhandled type %T", value)
	}
}
