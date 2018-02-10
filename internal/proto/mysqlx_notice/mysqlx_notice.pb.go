// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mysqlx_notice.proto

/*
Package mysqlx_notice is a generated protocol buffer package.

Notices

A notice

* is sent from the server to the client
* may be global or relate to the current message sequence

It is generated from these files:
	mysqlx_notice.proto

It has these top-level messages:
	Frame
	Warning
	SessionVariableChanged
	SessionStateChanged
*/
package mysqlx_notice

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/AlekSi/mysqlx/internal/proto/mysqlx"
import Mysqlx_Datatypes "github.com/AlekSi/mysqlx/internal/proto/mysqlx_datatypes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Frame_Scope int32

const (
	Frame_GLOBAL Frame_Scope = 1
	Frame_LOCAL  Frame_Scope = 2
)

var Frame_Scope_name = map[int32]string{
	1: "GLOBAL",
	2: "LOCAL",
}
var Frame_Scope_value = map[string]int32{
	"GLOBAL": 1,
	"LOCAL":  2,
}

func (x Frame_Scope) Enum() *Frame_Scope {
	p := new(Frame_Scope)
	*p = x
	return p
}
func (x Frame_Scope) String() string {
	return proto.EnumName(Frame_Scope_name, int32(x))
}
func (x *Frame_Scope) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Frame_Scope_value, data, "Frame_Scope")
	if err != nil {
		return err
	}
	*x = Frame_Scope(value)
	return nil
}
func (Frame_Scope) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Frame_Type int32

const (
	Frame_WARNING                  Frame_Type = 1
	Frame_SESSION_VARIABLE_CHANGED Frame_Type = 2
	Frame_SESSION_STATE_CHANGED    Frame_Type = 3
)

var Frame_Type_name = map[int32]string{
	1: "WARNING",
	2: "SESSION_VARIABLE_CHANGED",
	3: "SESSION_STATE_CHANGED",
}
var Frame_Type_value = map[string]int32{
	"WARNING":                  1,
	"SESSION_VARIABLE_CHANGED": 2,
	"SESSION_STATE_CHANGED":    3,
}

func (x Frame_Type) Enum() *Frame_Type {
	p := new(Frame_Type)
	*p = x
	return p
}
func (x Frame_Type) String() string {
	return proto.EnumName(Frame_Type_name, int32(x))
}
func (x *Frame_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Frame_Type_value, data, "Frame_Type")
	if err != nil {
		return err
	}
	*x = Frame_Type(value)
	return nil
}
func (Frame_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

type Warning_Level int32

const (
	Warning_NOTE    Warning_Level = 1
	Warning_WARNING Warning_Level = 2
	Warning_ERROR   Warning_Level = 3
)

var Warning_Level_name = map[int32]string{
	1: "NOTE",
	2: "WARNING",
	3: "ERROR",
}
var Warning_Level_value = map[string]int32{
	"NOTE":    1,
	"WARNING": 2,
	"ERROR":   3,
}

func (x Warning_Level) Enum() *Warning_Level {
	p := new(Warning_Level)
	*p = x
	return p
}
func (x Warning_Level) String() string {
	return proto.EnumName(Warning_Level_name, int32(x))
}
func (x *Warning_Level) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Warning_Level_value, data, "Warning_Level")
	if err != nil {
		return err
	}
	*x = Warning_Level(value)
	return nil
}
func (Warning_Level) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type SessionStateChanged_Parameter int32

const (
	SessionStateChanged_CURRENT_SCHEMA      SessionStateChanged_Parameter = 1
	SessionStateChanged_ACCOUNT_EXPIRED     SessionStateChanged_Parameter = 2
	SessionStateChanged_GENERATED_INSERT_ID SessionStateChanged_Parameter = 3
	SessionStateChanged_ROWS_AFFECTED       SessionStateChanged_Parameter = 4
	SessionStateChanged_ROWS_FOUND          SessionStateChanged_Parameter = 5
	SessionStateChanged_ROWS_MATCHED        SessionStateChanged_Parameter = 6
	SessionStateChanged_TRX_COMMITTED       SessionStateChanged_Parameter = 7
	SessionStateChanged_TRX_ROLLEDBACK      SessionStateChanged_Parameter = 9
	SessionStateChanged_PRODUCED_MESSAGE    SessionStateChanged_Parameter = 10
	SessionStateChanged_CLIENT_ID_ASSIGNED  SessionStateChanged_Parameter = 11
)

var SessionStateChanged_Parameter_name = map[int32]string{
	1:  "CURRENT_SCHEMA",
	2:  "ACCOUNT_EXPIRED",
	3:  "GENERATED_INSERT_ID",
	4:  "ROWS_AFFECTED",
	5:  "ROWS_FOUND",
	6:  "ROWS_MATCHED",
	7:  "TRX_COMMITTED",
	9:  "TRX_ROLLEDBACK",
	10: "PRODUCED_MESSAGE",
	11: "CLIENT_ID_ASSIGNED",
}
var SessionStateChanged_Parameter_value = map[string]int32{
	"CURRENT_SCHEMA":      1,
	"ACCOUNT_EXPIRED":     2,
	"GENERATED_INSERT_ID": 3,
	"ROWS_AFFECTED":       4,
	"ROWS_FOUND":          5,
	"ROWS_MATCHED":        6,
	"TRX_COMMITTED":       7,
	"TRX_ROLLEDBACK":      9,
	"PRODUCED_MESSAGE":    10,
	"CLIENT_ID_ASSIGNED":  11,
}

func (x SessionStateChanged_Parameter) Enum() *SessionStateChanged_Parameter {
	p := new(SessionStateChanged_Parameter)
	*p = x
	return p
}
func (x SessionStateChanged_Parameter) String() string {
	return proto.EnumName(SessionStateChanged_Parameter_name, int32(x))
}
func (x *SessionStateChanged_Parameter) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SessionStateChanged_Parameter_value, data, "SessionStateChanged_Parameter")
	if err != nil {
		return err
	}
	*x = SessionStateChanged_Parameter(value)
	return nil
}
func (SessionStateChanged_Parameter) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{3, 0}
}

// Common Frame for all Notices
//
// ===================================================== =====
// .type                                                 value
// ===================================================== =====
// :protobuf:msg:`Mysqlx.Notice::Warning`                1
// :protobuf:msg:`Mysqlx.Notice::SessionVariableChanged` 2
// :protobuf:msg:`Mysqlx.Notice::SessionStateChanged`    3
// ===================================================== =====
//
// :param type: the type of the payload
// :param payload: the payload of the notification
// :param scope: global or local notification
//
type Frame struct {
	Type             *uint32      `protobuf:"varint,1,req,name=type" json:"type,omitempty"`
	Scope            *Frame_Scope `protobuf:"varint,2,opt,name=scope,enum=Mysqlx.Notice.Frame_Scope,def=1" json:"scope,omitempty"`
	Payload          []byte       `protobuf:"bytes,3,opt,name=payload" json:"payload,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Frame) Reset()                    { *m = Frame{} }
func (m *Frame) String() string            { return proto.CompactTextString(m) }
func (*Frame) ProtoMessage()               {}
func (*Frame) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

const Default_Frame_Scope Frame_Scope = Frame_GLOBAL

func (m *Frame) GetType() uint32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *Frame) GetScope() Frame_Scope {
	if m != nil && m.Scope != nil {
		return *m.Scope
	}
	return Default_Frame_Scope
}

func (m *Frame) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

// Server-side warnings and notes
//
// ``.scope`` == ``local``
//   ``.level``, ``.code`` and ``.msg`` map the content of
//
//   .. code-block:: sql
//
//     SHOW WARNINGS
//
// ``.scope`` == ``global``
//   (undefined) will be used for global, unstructured messages like:
//
//   * server is shutting down
//   * a node disconnected from group
//   * schema or table dropped
//
// ========================================== =======================
// :protobuf:msg:`Mysqlx.Notice::Frame` field value
// ========================================== =======================
// ``.type``                                  1
// ``.scope``                                 ``local`` or ``global``
// ========================================== =======================
//
// :param level: warning level: Note or Warning
// :param code: warning code
// :param msg: warning message
type Warning struct {
	Level            *Warning_Level `protobuf:"varint,1,opt,name=level,enum=Mysqlx.Notice.Warning_Level,def=2" json:"level,omitempty"`
	Code             *uint32        `protobuf:"varint,2,req,name=code" json:"code,omitempty"`
	Msg              *string        `protobuf:"bytes,3,req,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *Warning) Reset()                    { *m = Warning{} }
func (m *Warning) String() string            { return proto.CompactTextString(m) }
func (*Warning) ProtoMessage()               {}
func (*Warning) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

const Default_Warning_Level Warning_Level = Warning_WARNING

func (m *Warning) GetLevel() Warning_Level {
	if m != nil && m.Level != nil {
		return *m.Level
	}
	return Default_Warning_Level
}

func (m *Warning) GetCode() uint32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *Warning) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

// Notify clients about changes to the current session variables
//
// Every change to a variable that is accessible through:
//
// .. code-block:: sql
//
//   SHOW SESSION VARIABLES
//
// ========================================== =========
// :protobuf:msg:`Mysqlx.Notice::Frame` field value
// ========================================== =========
// ``.type``                                  2
// ``.scope``                                 ``local``
// ========================================== =========
//
// :param namespace: namespace that param belongs to
// :param param: name of the variable
// :param value: the changed value of param
type SessionVariableChanged struct {
	Param            *string                  `protobuf:"bytes,1,req,name=param" json:"param,omitempty"`
	Value            *Mysqlx_Datatypes.Scalar `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *SessionVariableChanged) Reset()                    { *m = SessionVariableChanged{} }
func (m *SessionVariableChanged) String() string            { return proto.CompactTextString(m) }
func (*SessionVariableChanged) ProtoMessage()               {}
func (*SessionVariableChanged) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SessionVariableChanged) GetParam() string {
	if m != nil && m.Param != nil {
		return *m.Param
	}
	return ""
}

func (m *SessionVariableChanged) GetValue() *Mysqlx_Datatypes.Scalar {
	if m != nil {
		return m.Value
	}
	return nil
}

// Notify clients about changes to the internal session state
//
// ========================================== =========
// :protobuf:msg:`Mysqlx.Notice::Frame` field value
// ========================================== =========
// ``.type``                                  3
// ``.scope``                                 ``local``
// ========================================== =========
//
// :param param: parameter key
// :param value: updated value
type SessionStateChanged struct {
	Param            *SessionStateChanged_Parameter `protobuf:"varint,1,req,name=param,enum=Mysqlx.Notice.SessionStateChanged_Parameter" json:"param,omitempty"`
	Value            *Mysqlx_Datatypes.Scalar       `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte                         `json:"-"`
}

func (m *SessionStateChanged) Reset()                    { *m = SessionStateChanged{} }
func (m *SessionStateChanged) String() string            { return proto.CompactTextString(m) }
func (*SessionStateChanged) ProtoMessage()               {}
func (*SessionStateChanged) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SessionStateChanged) GetParam() SessionStateChanged_Parameter {
	if m != nil && m.Param != nil {
		return *m.Param
	}
	return SessionStateChanged_CURRENT_SCHEMA
}

func (m *SessionStateChanged) GetValue() *Mysqlx_Datatypes.Scalar {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*Frame)(nil), "Mysqlx.Notice.Frame")
	proto.RegisterType((*Warning)(nil), "Mysqlx.Notice.Warning")
	proto.RegisterType((*SessionVariableChanged)(nil), "Mysqlx.Notice.SessionVariableChanged")
	proto.RegisterType((*SessionStateChanged)(nil), "Mysqlx.Notice.SessionStateChanged")
	proto.RegisterEnum("Mysqlx.Notice.Frame_Scope", Frame_Scope_name, Frame_Scope_value)
	proto.RegisterEnum("Mysqlx.Notice.Frame_Type", Frame_Type_name, Frame_Type_value)
	proto.RegisterEnum("Mysqlx.Notice.Warning_Level", Warning_Level_name, Warning_Level_value)
	proto.RegisterEnum("Mysqlx.Notice.SessionStateChanged_Parameter", SessionStateChanged_Parameter_name, SessionStateChanged_Parameter_value)
}

func init() { proto.RegisterFile("mysqlx_notice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 595 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc7, 0x65, 0x27, 0x6e, 0xc8, 0xa4, 0x0d, 0xcb, 0xa6, 0xb4, 0x6e, 0x55, 0xa1, 0xc8, 0xa7,
	0x20, 0x21, 0x0b, 0xf5, 0x84, 0xca, 0x69, 0xb3, 0xde, 0xa4, 0x16, 0x8e, 0x5d, 0xed, 0x3a, 0x6d,
	0x4f, 0x58, 0xdb, 0xc4, 0x94, 0x20, 0x27, 0x0e, 0x8e, 0x5b, 0xd1, 0xb7, 0xe0, 0xc6, 0xa3, 0xf1,
	0x00, 0xbd, 0xf1, 0x14, 0x68, 0xd7, 0x49, 0x29, 0x15, 0x27, 0x6e, 0x3b, 0x33, 0xff, 0xdf, 0x7c,
	0x6a, 0xa1, 0x33, 0xbf, 0x5b, 0x7d, 0xcd, 0xbe, 0x25, 0x8b, 0xbc, 0x9c, 0x4d, 0x52, 0x77, 0x59,
	0xe4, 0x65, 0x8e, 0x77, 0x46, 0xda, 0xe9, 0x86, 0xda, 0x79, 0xb8, 0x5d, 0x69, 0xaa, 0xe0, 0xe1,
	0xde, 0x9a, 0x98, 0xca, 0x52, 0x96, 0x77, 0xcb, 0x74, 0x55, 0xf9, 0x9d, 0x7b, 0x03, 0xac, 0x41,
	0x21, 0xe7, 0x29, 0xc6, 0x50, 0x57, 0x01, 0xdb, 0xe8, 0x9a, 0xbd, 0x1d, 0xae, 0xdf, 0xf8, 0x1d,
	0x58, 0xab, 0x49, 0xbe, 0x4c, 0x6d, 0xb3, 0x6b, 0xf4, 0xda, 0xc7, 0x87, 0xee, 0x5f, 0x25, 0x5c,
	0x0d, 0xba, 0x42, 0x29, 0x4e, 0xb6, 0x86, 0x41, 0xd4, 0x27, 0x01, 0xaf, 0x00, 0x6c, 0x43, 0x63,
	0x29, 0xef, 0xb2, 0x5c, 0x4e, 0xed, 0x5a, 0xd7, 0xe8, 0x6d, 0xf3, 0x8d, 0xe9, 0xbc, 0x02, 0x4b,
	0x13, 0x18, 0x60, 0xcd, 0x20, 0x03, 0x37, 0xc1, 0x0a, 0x22, 0x4a, 0x02, 0x64, 0x3a, 0x01, 0xd4,
	0x63, 0x55, 0xbb, 0x05, 0x8d, 0x0b, 0xc2, 0x43, 0x3f, 0x1c, 0x22, 0x03, 0x1f, 0x81, 0x2d, 0x98,
	0x10, 0x7e, 0x14, 0x26, 0xe7, 0x84, 0xfb, 0xa4, 0x1f, 0xb0, 0x84, 0x9e, 0x92, 0x70, 0xc8, 0x3c,
	0x64, 0xe2, 0x03, 0x78, 0xb9, 0x89, 0x8a, 0x98, 0xc4, 0x7f, 0x42, 0xb5, 0x93, 0xfa, 0xf7, 0x5f,
	0x6f, 0x5b, 0xce, 0x0f, 0x03, 0x1a, 0x17, 0xb2, 0x58, 0xcc, 0x16, 0xd7, 0xf8, 0x3d, 0x58, 0x59,
	0x7a, 0x9b, 0x66, 0xb6, 0xa1, 0x67, 0x3a, 0x7a, 0x32, 0xd3, 0x5a, 0xe6, 0x06, 0x4a, 0x73, 0xb2,
	0x69, 0x81, 0x57, 0x8c, 0x5a, 0xd2, 0x24, 0x9f, 0xaa, 0x7d, 0xe8, 0x25, 0xa9, 0x37, 0x46, 0x50,
	0x9b, 0xaf, 0xae, 0xed, 0x5a, 0xd7, 0xec, 0x35, 0xb9, 0x7a, 0x3a, 0xaf, 0xc1, 0xd2, 0x38, 0x7e,
	0x06, 0xf5, 0x30, 0x8a, 0x19, 0x32, 0x1e, 0x4f, 0x63, 0xaa, 0x69, 0x19, 0xe7, 0x11, 0x47, 0x35,
	0xe7, 0x23, 0xec, 0x89, 0x74, 0xb5, 0x9a, 0xe5, 0x8b, 0x73, 0x59, 0xcc, 0xe4, 0x55, 0x96, 0xd2,
	0xcf, 0x72, 0x71, 0x9d, 0x4e, 0xf1, 0x2e, 0x58, 0x4b, 0x59, 0xc8, 0xb9, 0x3e, 0x48, 0x93, 0x57,
	0x06, 0x76, 0xc1, 0xba, 0x95, 0xd9, 0x4d, 0x75, 0x91, 0xd6, 0xb1, 0xbd, 0xe9, 0xde, 0x7b, 0xb8,
	0xab, 0x98, 0xc8, 0x4c, 0x16, 0xbc, 0x92, 0x39, 0xf7, 0x26, 0x74, 0xd6, 0x05, 0x44, 0x29, 0xcb,
	0x87, 0xec, 0xfd, 0xc7, 0xd9, 0xdb, 0xc7, 0x6f, 0x9e, 0x6c, 0xe1, 0x1f, 0x88, 0x7b, 0xa6, 0xf4,
	0x69, 0x99, 0x16, 0xff, 0xdb, 0xcb, 0x4f, 0x03, 0x9a, 0x0f, 0x49, 0x30, 0x86, 0x36, 0x1d, 0x73,
	0xce, 0xc2, 0x38, 0x11, 0xf4, 0x94, 0x8d, 0x08, 0x32, 0x70, 0x07, 0x9e, 0x13, 0x4a, 0xa3, 0x71,
	0x18, 0x27, 0xec, 0xf2, 0xcc, 0xe7, 0xfa, 0xba, 0xfb, 0xd0, 0x19, 0xb2, 0x90, 0x71, 0x12, 0x33,
	0x2f, 0xf1, 0x43, 0xc1, 0x78, 0x9c, 0xf8, 0x1e, 0xaa, 0xe1, 0x17, 0xb0, 0xc3, 0xa3, 0x0b, 0x91,
	0x90, 0xc1, 0x80, 0xd1, 0x98, 0x79, 0xa8, 0x8e, 0xdb, 0x00, 0xda, 0x35, 0x88, 0xc6, 0xa1, 0x87,
	0x2c, 0x8c, 0x60, 0x5b, 0xdb, 0x23, 0x12, 0xd3, 0x53, 0xe6, 0xa1, 0x2d, 0x05, 0xc5, 0xfc, 0x32,
	0xa1, 0xd1, 0x68, 0xe4, 0xc7, 0x0a, 0x6a, 0xa8, 0x4e, 0x94, 0x8b, 0x47, 0x41, 0xc0, 0xbc, 0x3e,
	0xa1, 0x1f, 0x50, 0x13, 0xef, 0x02, 0x3a, 0xe3, 0x91, 0x37, 0xa6, 0xcc, 0x4b, 0x46, 0x4c, 0x08,
	0x32, 0x64, 0x08, 0xf0, 0x1e, 0x60, 0x1a, 0xf8, 0xaa, 0x65, 0xdf, 0x4b, 0x88, 0x10, 0xfe, 0x30,
	0x64, 0x1e, 0x6a, 0xf5, 0x0f, 0x60, 0x7f, 0x92, 0xcf, 0x5d, 0xfd, 0xc7, 0xdc, 0xc9, 0x17, 0x77,
	0xfd, 0xeb, 0xae, 0x6e, 0x3e, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xef, 0xfe, 0x71, 0xf5, 0xab,
	0x03, 0x00, 0x00,
}
