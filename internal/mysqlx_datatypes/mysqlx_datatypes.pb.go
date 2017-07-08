// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mysqlx_datatypes.proto

/*
Package mysqlx_datatypes is a generated protocol buffer package.

It is generated from these files:
	mysqlx_datatypes.proto

It has these top-level messages:
	Scalar
	Object
	Array
	Any
*/
package mysqlx_datatypes

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Scalar_Type int32

const (
	Scalar_V_SINT   Scalar_Type = 1
	Scalar_V_UINT   Scalar_Type = 2
	Scalar_V_NULL   Scalar_Type = 3
	Scalar_V_OCTETS Scalar_Type = 4
	Scalar_V_DOUBLE Scalar_Type = 5
	Scalar_V_FLOAT  Scalar_Type = 6
	Scalar_V_BOOL   Scalar_Type = 7
	Scalar_V_STRING Scalar_Type = 8
)

var Scalar_Type_name = map[int32]string{
	1: "V_SINT",
	2: "V_UINT",
	3: "V_NULL",
	4: "V_OCTETS",
	5: "V_DOUBLE",
	6: "V_FLOAT",
	7: "V_BOOL",
	8: "V_STRING",
}
var Scalar_Type_value = map[string]int32{
	"V_SINT":   1,
	"V_UINT":   2,
	"V_NULL":   3,
	"V_OCTETS": 4,
	"V_DOUBLE": 5,
	"V_FLOAT":  6,
	"V_BOOL":   7,
	"V_STRING": 8,
}

func (x Scalar_Type) Enum() *Scalar_Type {
	p := new(Scalar_Type)
	*p = x
	return p
}
func (x Scalar_Type) String() string {
	return proto.EnumName(Scalar_Type_name, int32(x))
}
func (x *Scalar_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Scalar_Type_value, data, "Scalar_Type")
	if err != nil {
		return err
	}
	*x = Scalar_Type(value)
	return nil
}
func (Scalar_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Any_Type int32

const (
	Any_SCALAR Any_Type = 1
	Any_OBJECT Any_Type = 2
	Any_ARRAY  Any_Type = 3
)

var Any_Type_name = map[int32]string{
	1: "SCALAR",
	2: "OBJECT",
	3: "ARRAY",
}
var Any_Type_value = map[string]int32{
	"SCALAR": 1,
	"OBJECT": 2,
	"ARRAY":  3,
}

func (x Any_Type) Enum() *Any_Type {
	p := new(Any_Type)
	*p = x
	return p
}
func (x Any_Type) String() string {
	return proto.EnumName(Any_Type_name, int32(x))
}
func (x *Any_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Any_Type_value, data, "Any_Type")
	if err != nil {
		return err
	}
	*x = Any_Type(value)
	return nil
}
func (Any_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

// a scalar
type Scalar struct {
	Type         *Scalar_Type `protobuf:"varint,1,req,name=type,enum=Mysqlx.Datatypes.Scalar_Type" json:"type,omitempty"`
	VSignedInt   *int64       `protobuf:"zigzag64,2,opt,name=v_signed_int,json=vSignedInt" json:"v_signed_int,omitempty"`
	VUnsignedInt *uint64      `protobuf:"varint,3,opt,name=v_unsigned_int,json=vUnsignedInt" json:"v_unsigned_int,omitempty"`
	// 4 is unused, was Null which doesn't have a storage anymore
	VOctets          *Scalar_Octets `protobuf:"bytes,5,opt,name=v_octets,json=vOctets" json:"v_octets,omitempty"`
	VDouble          *float64       `protobuf:"fixed64,6,opt,name=v_double,json=vDouble" json:"v_double,omitempty"`
	VFloat           *float32       `protobuf:"fixed32,7,opt,name=v_float,json=vFloat" json:"v_float,omitempty"`
	VBool            *bool          `protobuf:"varint,8,opt,name=v_bool,json=vBool" json:"v_bool,omitempty"`
	VString          *Scalar_String `protobuf:"bytes,9,opt,name=v_string,json=vString" json:"v_string,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *Scalar) Reset()                    { *m = Scalar{} }
func (m *Scalar) String() string            { return proto.CompactTextString(m) }
func (*Scalar) ProtoMessage()               {}
func (*Scalar) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Scalar) GetType() Scalar_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Scalar_V_SINT
}

func (m *Scalar) GetVSignedInt() int64 {
	if m != nil && m.VSignedInt != nil {
		return *m.VSignedInt
	}
	return 0
}

func (m *Scalar) GetVUnsignedInt() uint64 {
	if m != nil && m.VUnsignedInt != nil {
		return *m.VUnsignedInt
	}
	return 0
}

func (m *Scalar) GetVOctets() *Scalar_Octets {
	if m != nil {
		return m.VOctets
	}
	return nil
}

func (m *Scalar) GetVDouble() float64 {
	if m != nil && m.VDouble != nil {
		return *m.VDouble
	}
	return 0
}

func (m *Scalar) GetVFloat() float32 {
	if m != nil && m.VFloat != nil {
		return *m.VFloat
	}
	return 0
}

func (m *Scalar) GetVBool() bool {
	if m != nil && m.VBool != nil {
		return *m.VBool
	}
	return false
}

func (m *Scalar) GetVString() *Scalar_String {
	if m != nil {
		return m.VString
	}
	return nil
}

// a string with a charset/collation
type Scalar_String struct {
	Value            []byte  `protobuf:"bytes,1,req,name=value" json:"value,omitempty"`
	Collation        *uint64 `protobuf:"varint,2,opt,name=collation" json:"collation,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Scalar_String) Reset()                    { *m = Scalar_String{} }
func (m *Scalar_String) String() string            { return proto.CompactTextString(m) }
func (*Scalar_String) ProtoMessage()               {}
func (*Scalar_String) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Scalar_String) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Scalar_String) GetCollation() uint64 {
	if m != nil && m.Collation != nil {
		return *m.Collation
	}
	return 0
}

// an opaque octet sequence, with an optional content_type
// See ``Mysqlx.Resultset.ColumnMetadata`` for list of known values.
type Scalar_Octets struct {
	Value            []byte  `protobuf:"bytes,1,req,name=value" json:"value,omitempty"`
	ContentType      *uint32 `protobuf:"varint,2,opt,name=content_type,json=contentType" json:"content_type,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Scalar_Octets) Reset()                    { *m = Scalar_Octets{} }
func (m *Scalar_Octets) String() string            { return proto.CompactTextString(m) }
func (*Scalar_Octets) ProtoMessage()               {}
func (*Scalar_Octets) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *Scalar_Octets) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Scalar_Octets) GetContentType() uint32 {
	if m != nil && m.ContentType != nil {
		return *m.ContentType
	}
	return 0
}

// a object
type Object struct {
	Fld              []*Object_ObjectField `protobuf:"bytes,1,rep,name=fld" json:"fld,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *Object) Reset()                    { *m = Object{} }
func (m *Object) String() string            { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()               {}
func (*Object) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Object) GetFld() []*Object_ObjectField {
	if m != nil {
		return m.Fld
	}
	return nil
}

type Object_ObjectField struct {
	Key              *string `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            *Any    `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Object_ObjectField) Reset()                    { *m = Object_ObjectField{} }
func (m *Object_ObjectField) String() string            { return proto.CompactTextString(m) }
func (*Object_ObjectField) ProtoMessage()               {}
func (*Object_ObjectField) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *Object_ObjectField) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Object_ObjectField) GetValue() *Any {
	if m != nil {
		return m.Value
	}
	return nil
}

// a Array
type Array struct {
	Value            []*Any `protobuf:"bytes,1,rep,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Array) Reset()                    { *m = Array{} }
func (m *Array) String() string            { return proto.CompactTextString(m) }
func (*Array) ProtoMessage()               {}
func (*Array) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Array) GetValue() []*Any {
	if m != nil {
		return m.Value
	}
	return nil
}

// a helper to allow all field types
type Any struct {
	Type             *Any_Type `protobuf:"varint,1,req,name=type,enum=Mysqlx.Datatypes.Any_Type" json:"type,omitempty"`
	Scalar           *Scalar   `protobuf:"bytes,2,opt,name=scalar" json:"scalar,omitempty"`
	Obj              *Object   `protobuf:"bytes,3,opt,name=obj" json:"obj,omitempty"`
	Array            *Array    `protobuf:"bytes,4,opt,name=array" json:"array,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *Any) Reset()                    { *m = Any{} }
func (m *Any) String() string            { return proto.CompactTextString(m) }
func (*Any) ProtoMessage()               {}
func (*Any) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Any) GetType() Any_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Any_SCALAR
}

func (m *Any) GetScalar() *Scalar {
	if m != nil {
		return m.Scalar
	}
	return nil
}

func (m *Any) GetObj() *Object {
	if m != nil {
		return m.Obj
	}
	return nil
}

func (m *Any) GetArray() *Array {
	if m != nil {
		return m.Array
	}
	return nil
}

func init() {
	proto.RegisterType((*Scalar)(nil), "Mysqlx.Datatypes.Scalar")
	proto.RegisterType((*Scalar_String)(nil), "Mysqlx.Datatypes.Scalar.String")
	proto.RegisterType((*Scalar_Octets)(nil), "Mysqlx.Datatypes.Scalar.Octets")
	proto.RegisterType((*Object)(nil), "Mysqlx.Datatypes.Object")
	proto.RegisterType((*Object_ObjectField)(nil), "Mysqlx.Datatypes.Object.ObjectField")
	proto.RegisterType((*Array)(nil), "Mysqlx.Datatypes.Array")
	proto.RegisterType((*Any)(nil), "Mysqlx.Datatypes.Any")
	proto.RegisterEnum("Mysqlx.Datatypes.Scalar_Type", Scalar_Type_name, Scalar_Type_value)
	proto.RegisterEnum("Mysqlx.Datatypes.Any_Type", Any_Type_name, Any_Type_value)
}

func init() { proto.RegisterFile("mysqlx_datatypes.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 585 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0xe5, 0xe6, 0xab, 0x3d, 0x29, 0x53, 0x64, 0x31, 0x16, 0xaa, 0x01, 0xa1, 0xda, 0x45,
	0x00, 0x11, 0x41, 0x85, 0xb8, 0x40, 0xdc, 0xa4, 0xfb, 0x40, 0x43, 0x61, 0x91, 0x9c, 0x76, 0x12,
	0x57, 0x51, 0x9a, 0x66, 0x53, 0x47, 0x16, 0x8f, 0xc6, 0xb5, 0x96, 0x97, 0xe0, 0x71, 0x78, 0x2b,
	0xde, 0x01, 0xd9, 0xce, 0xb4, 0xc2, 0x3a, 0xb8, 0xca, 0xff, 0xfc, 0xfb, 0x3b, 0xee, 0xf1, 0x39,
	0x47, 0x86, 0x47, 0x97, 0x4d, 0xfd, 0xbd, 0xbc, 0x4e, 0xe7, 0x19, 0xcb, 0x58, 0x73, 0x55, 0xd4,
	0xc1, 0xd5, 0x92, 0x32, 0x8a, 0x9d, 0x2f, 0xd2, 0x0f, 0x0e, 0x6e, 0xfc, 0xe1, 0x4f, 0x1d, 0xcc,
	0x24, 0xcf, 0xca, 0x6c, 0x89, 0xdf, 0x82, 0x2e, 0x3c, 0x17, 0x79, 0x1d, 0x7f, 0x6b, 0xf4, 0x24,
	0xf8, 0x9b, 0x0d, 0x14, 0x17, 0x4c, 0x9a, 0xab, 0x82, 0x48, 0x14, 0x7b, 0xd0, 0xe7, 0x69, 0xbd,
	0x38, 0xaf, 0x8a, 0x79, 0xba, 0xa8, 0x98, 0xdb, 0xf1, 0x90, 0x8f, 0x09, 0xf0, 0x44, 0x5a, 0xc7,
	0x15, 0xc3, 0x7b, 0xb0, 0xc5, 0xd3, 0x55, 0xb5, 0xc6, 0x68, 0x1e, 0xf2, 0x75, 0xd2, 0xe7, 0xd3,
	0xd6, 0x14, 0xd4, 0x07, 0xe8, 0xf2, 0x94, 0xe6, 0xac, 0x60, 0xb5, 0x6b, 0x78, 0xc8, 0xb7, 0x47,
	0xcf, 0xee, 0xfd, 0xfb, 0x58, 0x62, 0xc4, 0xe2, 0x4a, 0xe0, 0xc7, 0x22, 0x77, 0x4e, 0x57, 0xb3,
	0xb2, 0x70, 0x4d, 0x0f, 0xf9, 0x88, 0x58, 0xfc, 0x40, 0x86, 0x78, 0x07, 0x2c, 0x9e, 0x9e, 0x95,
	0x34, 0x63, 0xae, 0xe5, 0x21, 0xbf, 0x43, 0x4c, 0x7e, 0x24, 0x22, 0xbc, 0x0d, 0x26, 0x4f, 0x67,
	0x94, 0x96, 0x6e, 0xd7, 0x43, 0x7e, 0x97, 0x18, 0x7c, 0x4c, 0x69, 0xa9, 0xca, 0xa8, 0xd9, 0x72,
	0x51, 0x9d, 0xbb, 0xbd, 0xff, 0x94, 0x91, 0x48, 0x8c, 0x58, 0x5c, 0x89, 0xc1, 0x47, 0x30, 0x95,
	0xc2, 0x0f, 0xc1, 0xe0, 0x59, 0xb9, 0x52, 0x8d, 0xec, 0x13, 0x15, 0xe0, 0x5d, 0xe8, 0xe5, 0xb4,
	0x2c, 0x33, 0xb6, 0xa0, 0x95, 0xec, 0x93, 0x4e, 0x6e, 0x8d, 0x41, 0x08, 0x66, 0x7b, 0x9d, 0xcd,
	0xd9, 0xcf, 0xa1, 0x9f, 0xd3, 0x8a, 0x15, 0x15, 0x4b, 0xe5, 0x8c, 0xc4, 0x01, 0x0f, 0x88, 0xdd,
	0x7a, 0x62, 0x22, 0xc3, 0x4b, 0xd0, 0xc5, 0x17, 0x03, 0x98, 0xa7, 0x69, 0x72, 0x7c, 0x32, 0x71,
	0x90, 0xd2, 0x53, 0xa1, 0x3b, 0x4a, 0x9f, 0x4c, 0xa3, 0xc8, 0xd1, 0x70, 0x1f, 0xba, 0xa7, 0x69,
	0xbc, 0x3f, 0x39, 0x9c, 0x24, 0x8e, 0xae, 0xa2, 0x83, 0x78, 0x3a, 0x8e, 0x0e, 0x1d, 0x03, 0xdb,
	0x60, 0x9d, 0xa6, 0x47, 0x51, 0x1c, 0x4e, 0x1c, 0x53, 0x25, 0x8d, 0xe3, 0x38, 0x72, 0x2c, 0x85,
	0x25, 0x13, 0x72, 0x7c, 0xf2, 0xc9, 0xe9, 0x0e, 0x7f, 0x20, 0x30, 0xe3, 0xd9, 0x45, 0x91, 0x33,
	0xfc, 0x1e, 0xb4, 0xb3, 0x72, 0xee, 0x22, 0x4f, 0xf3, 0xed, 0xd1, 0xde, 0xdd, 0x8e, 0x29, 0xac,
	0xfd, 0x1c, 0x2d, 0x8a, 0x72, 0x4e, 0x44, 0xc2, 0x20, 0x02, 0x7b, 0xcd, 0xc3, 0x0e, 0x68, 0xdf,
	0x8a, 0x46, 0xde, 0xbb, 0x47, 0x84, 0xc4, 0xaf, 0x6e, 0x7a, 0xd1, 0xf1, 0x3a, 0xbe, 0x3d, 0xda,
	0xbe, 0x7b, 0x74, 0x58, 0x35, 0x6d, 0x8b, 0x86, 0xef, 0xc0, 0x08, 0x97, 0xcb, 0x6c, 0x2d, 0x4b,
	0x15, 0xf4, 0xef, 0xac, 0x5f, 0x08, 0xb4, 0xb0, 0x6a, 0x70, 0xf0, 0xc7, 0xf2, 0x0f, 0x36, 0xe6,
	0xac, 0x6f, 0xfe, 0x1b, 0x30, 0x6b, 0xb9, 0x08, 0x72, 0x14, 0xf6, 0xc8, 0xbd, 0x6f, 0x51, 0x48,
	0xcb, 0xe1, 0x97, 0xa0, 0xd1, 0xd9, 0x85, 0x5c, 0xff, 0x8d, 0xb8, 0x6a, 0x05, 0x11, 0x10, 0x7e,
	0x0d, 0x46, 0x26, 0xee, 0xe2, 0xea, 0x92, 0xde, 0xd9, 0x50, 0x8e, 0xf8, 0x99, 0x28, 0x6a, 0xf8,
	0xe2, 0x76, 0xf4, 0xc9, 0x7e, 0x18, 0x85, 0x44, 0x8d, 0x3e, 0x1e, 0x7f, 0x3e, 0xdc, 0x17, 0xa3,
	0xef, 0x81, 0x11, 0x12, 0x12, 0x7e, 0x75, 0xb4, 0xf1, 0x53, 0xd8, 0xcd, 0xe9, 0x65, 0x20, 0xdf,
	0x87, 0x20, 0xbf, 0x50, 0xe2, 0x5a, 0x3d, 0x0f, 0xb3, 0xd5, 0xd9, 0xef, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xc2, 0x2e, 0x09, 0x9b, 0x3a, 0x04, 0x00, 0x00,
}
