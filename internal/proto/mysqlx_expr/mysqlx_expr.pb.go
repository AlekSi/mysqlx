// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mysqlx_expr.proto

/*
Package mysqlx_expr is a generated protocol buffer package.

Expression syntax

expr is the fundamental structure in various places
of the SQL language:

* ``SELECT <expr> AS ...``
* ``WHERE <expr>``

The structures can be used to:

* build an Item-tree in the MySQL Server
* generate SQL from it
* use as filter condition in CRUD's Find(), Update() and Delete() calls.

It is generated from these files:
	mysqlx_expr.proto

It has these top-level messages:
	Expr
	Identifier
	DocumentPathItem
	ColumnIdentifier
	FunctionCall
	Operator
	Object
	Array
*/
package mysqlx_expr

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
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

type Expr_Type int32

const (
	Expr_IDENT       Expr_Type = 1
	Expr_LITERAL     Expr_Type = 2
	Expr_VARIABLE    Expr_Type = 3
	Expr_FUNC_CALL   Expr_Type = 4
	Expr_OPERATOR    Expr_Type = 5
	Expr_PLACEHOLDER Expr_Type = 6
	Expr_OBJECT      Expr_Type = 7
	Expr_ARRAY       Expr_Type = 8
)

var Expr_Type_name = map[int32]string{
	1: "IDENT",
	2: "LITERAL",
	3: "VARIABLE",
	4: "FUNC_CALL",
	5: "OPERATOR",
	6: "PLACEHOLDER",
	7: "OBJECT",
	8: "ARRAY",
}
var Expr_Type_value = map[string]int32{
	"IDENT":       1,
	"LITERAL":     2,
	"VARIABLE":    3,
	"FUNC_CALL":   4,
	"OPERATOR":    5,
	"PLACEHOLDER": 6,
	"OBJECT":      7,
	"ARRAY":       8,
}

func (x Expr_Type) Enum() *Expr_Type {
	p := new(Expr_Type)
	*p = x
	return p
}
func (x Expr_Type) String() string {
	return proto.EnumName(Expr_Type_name, int32(x))
}
func (x *Expr_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Expr_Type_value, data, "Expr_Type")
	if err != nil {
		return err
	}
	*x = Expr_Type(value)
	return nil
}
func (Expr_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type DocumentPathItem_Type int32

const (
	DocumentPathItem_MEMBER               DocumentPathItem_Type = 1
	DocumentPathItem_MEMBER_ASTERISK      DocumentPathItem_Type = 2
	DocumentPathItem_ARRAY_INDEX          DocumentPathItem_Type = 3
	DocumentPathItem_ARRAY_INDEX_ASTERISK DocumentPathItem_Type = 4
	DocumentPathItem_DOUBLE_ASTERISK      DocumentPathItem_Type = 5
)

var DocumentPathItem_Type_name = map[int32]string{
	1: "MEMBER",
	2: "MEMBER_ASTERISK",
	3: "ARRAY_INDEX",
	4: "ARRAY_INDEX_ASTERISK",
	5: "DOUBLE_ASTERISK",
}
var DocumentPathItem_Type_value = map[string]int32{
	"MEMBER":               1,
	"MEMBER_ASTERISK":      2,
	"ARRAY_INDEX":          3,
	"ARRAY_INDEX_ASTERISK": 4,
	"DOUBLE_ASTERISK":      5,
}

func (x DocumentPathItem_Type) Enum() *DocumentPathItem_Type {
	p := new(DocumentPathItem_Type)
	*p = x
	return p
}
func (x DocumentPathItem_Type) String() string {
	return proto.EnumName(DocumentPathItem_Type_name, int32(x))
}
func (x *DocumentPathItem_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(DocumentPathItem_Type_value, data, "DocumentPathItem_Type")
	if err != nil {
		return err
	}
	*x = DocumentPathItem_Type(value)
	return nil
}
func (DocumentPathItem_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

// Expressions
//
// the "root" of the expression tree
//
// .. productionlist::
//   expr: `operator` |
//       : `identifier` |
//       : `function_call` |
//       : variable |
//       : `literal` |
//       : placeholder
//
// If expression type is PLACEHOLDER then it refers to the value of a parameter
// specified when executing a statement (see `args` field of `StmtExecute` command).
// Field `position` (which must be present for such an expression) gives 0-based
// position of the parameter in the parameter list.
//
type Expr struct {
	Type             *Expr_Type               `protobuf:"varint,1,req,name=type,enum=Mysqlx.Expr.Expr_Type" json:"type,omitempty"`
	Identifier       *ColumnIdentifier        `protobuf:"bytes,2,opt,name=identifier" json:"identifier,omitempty"`
	Variable         *string                  `protobuf:"bytes,3,opt,name=variable" json:"variable,omitempty"`
	Literal          *Mysqlx_Datatypes.Scalar `protobuf:"bytes,4,opt,name=literal" json:"literal,omitempty"`
	FunctionCall     *FunctionCall            `protobuf:"bytes,5,opt,name=function_call,json=functionCall" json:"function_call,omitempty"`
	Operator         *Operator                `protobuf:"bytes,6,opt,name=operator" json:"operator,omitempty"`
	Position         *uint32                  `protobuf:"varint,7,opt,name=position" json:"position,omitempty"`
	Object           *Object                  `protobuf:"bytes,8,opt,name=object" json:"object,omitempty"`
	Array            *Array                   `protobuf:"bytes,9,opt,name=array" json:"array,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *Expr) Reset()                    { *m = Expr{} }
func (m *Expr) String() string            { return proto.CompactTextString(m) }
func (*Expr) ProtoMessage()               {}
func (*Expr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Expr) GetType() Expr_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Expr_IDENT
}

func (m *Expr) GetIdentifier() *ColumnIdentifier {
	if m != nil {
		return m.Identifier
	}
	return nil
}

func (m *Expr) GetVariable() string {
	if m != nil && m.Variable != nil {
		return *m.Variable
	}
	return ""
}

func (m *Expr) GetLiteral() *Mysqlx_Datatypes.Scalar {
	if m != nil {
		return m.Literal
	}
	return nil
}

func (m *Expr) GetFunctionCall() *FunctionCall {
	if m != nil {
		return m.FunctionCall
	}
	return nil
}

func (m *Expr) GetOperator() *Operator {
	if m != nil {
		return m.Operator
	}
	return nil
}

func (m *Expr) GetPosition() uint32 {
	if m != nil && m.Position != nil {
		return *m.Position
	}
	return 0
}

func (m *Expr) GetObject() *Object {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *Expr) GetArray() *Array {
	if m != nil {
		return m.Array
	}
	return nil
}

// identifier: name, schame.name
//
// .. productionlist::
//   identifier: string "." string |
//             : string
type Identifier struct {
	Name             *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	SchemaName       *string `protobuf:"bytes,2,opt,name=schema_name,json=schemaName" json:"schema_name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Identifier) Reset()                    { *m = Identifier{} }
func (m *Identifier) String() string            { return proto.CompactTextString(m) }
func (*Identifier) ProtoMessage()               {}
func (*Identifier) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Identifier) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Identifier) GetSchemaName() string {
	if m != nil && m.SchemaName != nil {
		return *m.SchemaName
	}
	return ""
}

// DocumentPathItem
//
// .. productionlist::
//    document_path: path_item | path_item document_path
//    path_item    : member | array_index | "**"
//    member       : "." string | "." "*"
//    array_index  : "[" number "]" | "[" "*" "]"
//
type DocumentPathItem struct {
	Type             *DocumentPathItem_Type `protobuf:"varint,1,req,name=type,enum=Mysqlx.Expr.DocumentPathItem_Type" json:"type,omitempty"`
	Value            *string                `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Index            *uint32                `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *DocumentPathItem) Reset()                    { *m = DocumentPathItem{} }
func (m *DocumentPathItem) String() string            { return proto.CompactTextString(m) }
func (*DocumentPathItem) ProtoMessage()               {}
func (*DocumentPathItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DocumentPathItem) GetType() DocumentPathItem_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return DocumentPathItem_MEMBER
}

func (m *DocumentPathItem) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func (m *DocumentPathItem) GetIndex() uint32 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

// col_identifier (table): col@doc_path, tbl.col@doc_path col, tbl.col, schema.tbl.col
// col_identifier (document): doc_path
//
// .. productionlist::
//   col_identifier: string "." string "." string |
//             : string "." string |
//             : string |
//             : string "." string "." string "@" document_path |
//             : string "." string "@" document_path |
//             : string "@" document_path |
//             : document_path
//    document_path: member | arrayLocation | doubleAsterisk
//    member = "." string | "." "*"
//    arrayLocation = "[" index "]" | "[" "*" "]"
//    doubleAsterisk = "**"
//
type ColumnIdentifier struct {
	DocumentPath     []*DocumentPathItem `protobuf:"bytes,1,rep,name=document_path,json=documentPath" json:"document_path,omitempty"`
	Name             *string             `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	TableName        *string             `protobuf:"bytes,3,opt,name=table_name,json=tableName" json:"table_name,omitempty"`
	SchemaName       *string             `protobuf:"bytes,4,opt,name=schema_name,json=schemaName" json:"schema_name,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *ColumnIdentifier) Reset()                    { *m = ColumnIdentifier{} }
func (m *ColumnIdentifier) String() string            { return proto.CompactTextString(m) }
func (*ColumnIdentifier) ProtoMessage()               {}
func (*ColumnIdentifier) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ColumnIdentifier) GetDocumentPath() []*DocumentPathItem {
	if m != nil {
		return m.DocumentPath
	}
	return nil
}

func (m *ColumnIdentifier) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ColumnIdentifier) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *ColumnIdentifier) GetSchemaName() string {
	if m != nil && m.SchemaName != nil {
		return *m.SchemaName
	}
	return ""
}

// function call: ``func(a, b, "1", 3)``
//
// .. productionlist::
//   function_call: `identifier` "(" [ `expr` ["," `expr` ]* ] ")"
type FunctionCall struct {
	Name             *Identifier `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Param            []*Expr     `protobuf:"bytes,2,rep,name=param" json:"param,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *FunctionCall) Reset()                    { *m = FunctionCall{} }
func (m *FunctionCall) String() string            { return proto.CompactTextString(m) }
func (*FunctionCall) ProtoMessage()               {}
func (*FunctionCall) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *FunctionCall) GetName() *Identifier {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *FunctionCall) GetParam() []*Expr {
	if m != nil {
		return m.Param
	}
	return nil
}

// operator: ``<<(a, b)``
//
// .. note::
//
//   Non-authoritative list of operators implemented (case sensitive):
//
//   Nullary
//     * ``*``
//     * ``default``
//
//   Unary
//     * ``!``
//     * ``sign_plus``
//     * ``sign_minus``
//     * ``~``
//
//   Binary
//     * ``&&``
//     * ``||``
//     * ``xor``
//     * ``==``
//     * ``!=``
//     * ``>``
//     * ``>=``
//     * ``<``
//     * ``<=``
//     * ``&``
//     * ``|``
//     * ``^``
//     * ``<<``
//     * ``>>``
//     * ``+``
//     * ``-``
//     * ``*``
//     * ``/``
//     * ``div``
//     * ``%``
//     * ``is``
//     * ``is_not``
//     * ``regexp``
//     * ``not_regexp``
//     * ``like``
//     * ``not_like``
//     * ``cast``
//     * ``cont_in``
//     * ``not_cont_in``
//
//   Using special representation, with more than 2 params
//     * ``in`` (param[0] IN (param[1], param[2], ...))
//     * ``not_in`` (param[0] NOT IN (param[1], param[2], ...))
//
//   Ternary
//     * ``between``
//     * ``between_not``
//     * ``date_add``
//     * ``date_sub``
//
//   Units for date_add/date_sub
//     * ``MICROSECOND``
//     * ``SECOND``
//     * ``MINUTE``
//     * ``HOUR``
//     * ``DAY``
//     * ``WEEK``
//     * ``MONTH``
//     * ``QUARTER``
//     * ``YEAR``
//     * ``SECOND_MICROSECOND``
//     * ``MINUTE_MICROSECOND``
//     * ``MINUTE_SECOND``
//     * ``HOUR_MICROSECOND``
//     * ``HOUR_SECOND``
//     * ``HOUR_MINUTE``
//     * ``DAY_MICROSECOND``
//     * ``DAY_SECOND``
//     * ``DAY_MINUTE``
//     * ``DAY_HOUR``
//
//   Types for cast
//     * ``BINARY[(N)]``
//     * ``CHAR[(N)]``
//     * ``DATE``
//     * ``DATETIME``
//     * ``DECIMAL[(M[,D])]``
//     * ``JSON``
//     * ``SIGNED [INTEGER]``
//     * ``TIME``
//     * ``UNSIGNED [INTEGER]``
//
// .. productionlist::
//   operator: `name` "(" [ `expr` ["," `expr` ]* ] ")"
type Operator struct {
	Name             *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Param            []*Expr `protobuf:"bytes,2,rep,name=param" json:"param,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Operator) Reset()                    { *m = Operator{} }
func (m *Operator) String() string            { return proto.CompactTextString(m) }
func (*Operator) ProtoMessage()               {}
func (*Operator) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Operator) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Operator) GetParam() []*Expr {
	if m != nil {
		return m.Param
	}
	return nil
}

// an object (with expression values)
type Object struct {
	Fld              []*Object_ObjectField `protobuf:"bytes,1,rep,name=fld" json:"fld,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *Object) Reset()                    { *m = Object{} }
func (m *Object) String() string            { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()               {}
func (*Object) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Object) GetFld() []*Object_ObjectField {
	if m != nil {
		return m.Fld
	}
	return nil
}

type Object_ObjectField struct {
	Key              *string `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            *Expr   `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Object_ObjectField) Reset()                    { *m = Object_ObjectField{} }
func (m *Object_ObjectField) String() string            { return proto.CompactTextString(m) }
func (*Object_ObjectField) ProtoMessage()               {}
func (*Object_ObjectField) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6, 0} }

func (m *Object_ObjectField) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Object_ObjectField) GetValue() *Expr {
	if m != nil {
		return m.Value
	}
	return nil
}

// a Array of expressions
type Array struct {
	Value            []*Expr `protobuf:"bytes,1,rep,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Array) Reset()                    { *m = Array{} }
func (m *Array) String() string            { return proto.CompactTextString(m) }
func (*Array) ProtoMessage()               {}
func (*Array) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Array) GetValue() []*Expr {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*Expr)(nil), "Mysqlx.Expr.Expr")
	proto.RegisterType((*Identifier)(nil), "Mysqlx.Expr.Identifier")
	proto.RegisterType((*DocumentPathItem)(nil), "Mysqlx.Expr.DocumentPathItem")
	proto.RegisterType((*ColumnIdentifier)(nil), "Mysqlx.Expr.ColumnIdentifier")
	proto.RegisterType((*FunctionCall)(nil), "Mysqlx.Expr.FunctionCall")
	proto.RegisterType((*Operator)(nil), "Mysqlx.Expr.Operator")
	proto.RegisterType((*Object)(nil), "Mysqlx.Expr.Object")
	proto.RegisterType((*Object_ObjectField)(nil), "Mysqlx.Expr.Object.ObjectField")
	proto.RegisterType((*Array)(nil), "Mysqlx.Expr.Array")
	proto.RegisterEnum("Mysqlx.Expr.Expr_Type", Expr_Type_name, Expr_Type_value)
	proto.RegisterEnum("Mysqlx.Expr.DocumentPathItem_Type", DocumentPathItem_Type_name, DocumentPathItem_Type_value)
}

func init() { proto.RegisterFile("mysqlx_expr.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 714 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x5b, 0x6e, 0xda, 0x40,
	0x14, 0x95, 0xc1, 0xe6, 0x71, 0x81, 0xc6, 0x99, 0xa4, 0x89, 0x13, 0x29, 0x0a, 0xf2, 0x4f, 0x51,
	0x23, 0xa1, 0x86, 0x8f, 0xfe, 0xb5, 0x92, 0x01, 0xa7, 0xa1, 0x75, 0x20, 0x9a, 0x90, 0xaa, 0xfd,
	0x42, 0x13, 0x7b, 0x68, 0x9c, 0xfa, 0x15, 0x63, 0x22, 0x58, 0x40, 0x37, 0xd1, 0x25, 0x74, 0x53,
	0xdd, 0x4a, 0x35, 0x33, 0x06, 0xec, 0x3c, 0xa4, 0xfe, 0xc0, 0xcc, 0x3d, 0xe7, 0x5e, 0xdf, 0xd7,
	0x19, 0xd8, 0xf6, 0x97, 0xb3, 0x7b, 0x6f, 0x31, 0xa1, 0x8b, 0x28, 0x6e, 0x47, 0x71, 0x98, 0x84,
	0xa8, 0x76, 0xc1, 0x4d, 0x6d, 0x73, 0x11, 0xc5, 0x87, 0x7b, 0x29, 0xee, 0x90, 0x84, 0x24, 0xcb,
	0x88, 0xce, 0x04, 0x49, 0xff, 0x2d, 0x83, 0xcc, 0x08, 0xe8, 0x2d, 0xc8, 0xcc, 0xae, 0x49, 0xcd,
	0x42, 0xeb, 0x55, 0x67, 0xaf, 0x9d, 0x71, 0x16, 0x3f, 0xe3, 0x65, 0x44, 0x31, 0xe7, 0xa0, 0x0f,
	0x00, 0xae, 0x43, 0x83, 0xc4, 0x9d, 0xba, 0x34, 0xd6, 0x0a, 0x4d, 0xa9, 0x55, 0xeb, 0x1c, 0xe5,
	0x3c, 0x7a, 0xa1, 0x37, 0xf7, 0x83, 0xc1, 0x9a, 0x84, 0x33, 0x0e, 0xe8, 0x10, 0x2a, 0x0f, 0x24,
	0x76, 0xc9, 0x8d, 0x47, 0xb5, 0x62, 0x53, 0x6a, 0x55, 0xf1, 0xfa, 0x8e, 0x3a, 0x50, 0xf6, 0xdc,
	0x84, 0xc6, 0xc4, 0xd3, 0x64, 0x1e, 0x57, 0x5b, 0xc5, 0xed, 0xaf, 0x33, 0xbf, 0xb2, 0x89, 0x47,
	0x62, 0xbc, 0x22, 0xa2, 0x8f, 0xd0, 0x98, 0xce, 0x03, 0x3b, 0x71, 0xc3, 0x60, 0x62, 0x13, 0xcf,
	0xd3, 0x14, 0xee, 0x79, 0x90, 0xcb, 0xe8, 0x2c, 0x65, 0xf4, 0x88, 0xe7, 0xe1, 0xfa, 0x34, 0x73,
	0x43, 0xa7, 0x50, 0x09, 0x23, 0x1a, 0x93, 0x24, 0x8c, 0xb5, 0x12, 0x77, 0x7d, 0x9d, 0x73, 0x1d,
	0xa5, 0x20, 0x5e, 0xd3, 0x58, 0x09, 0x51, 0x38, 0x73, 0x59, 0x08, 0xad, 0xdc, 0x94, 0x5a, 0x0d,
	0xbc, 0xbe, 0xa3, 0x13, 0x28, 0x85, 0x37, 0x77, 0xd4, 0x4e, 0xb4, 0x0a, 0x0f, 0xb6, 0x93, 0x0f,
	0xc6, 0x21, 0x9c, 0x52, 0x50, 0x0b, 0x14, 0x12, 0xc7, 0x64, 0xa9, 0x55, 0x39, 0x17, 0xe5, 0xb8,
	0x06, 0x43, 0xb0, 0x20, 0xe8, 0xf7, 0x20, 0xb3, 0x11, 0xa0, 0x2a, 0x28, 0x83, 0xbe, 0x39, 0x1c,
	0xab, 0x12, 0xaa, 0x41, 0xd9, 0x1a, 0x8c, 0x4d, 0x6c, 0x58, 0x6a, 0x01, 0xd5, 0xa1, 0xf2, 0xd5,
	0xc0, 0x03, 0xa3, 0x6b, 0x99, 0x6a, 0x11, 0x35, 0xa0, 0x7a, 0x76, 0x3d, 0xec, 0x4d, 0x7a, 0x86,
	0x65, 0xa9, 0x32, 0x03, 0x47, 0x97, 0x26, 0x36, 0xc6, 0x23, 0xac, 0x2a, 0x68, 0x0b, 0x6a, 0x97,
	0x96, 0xd1, 0x33, 0xcf, 0x47, 0x56, 0xdf, 0xc4, 0x6a, 0x09, 0x01, 0x94, 0x46, 0xdd, 0xcf, 0x66,
	0x6f, 0xac, 0x96, 0x59, 0x7c, 0x03, 0x63, 0xe3, 0xbb, 0x5a, 0xd1, 0x0d, 0x80, 0xcd, 0x08, 0x11,
	0x02, 0x39, 0x20, 0xbe, 0xd8, 0x90, 0x2a, 0xe6, 0x67, 0x74, 0x0c, 0xb5, 0x99, 0x7d, 0x4b, 0x7d,
	0x32, 0xe1, 0x50, 0x81, 0x4f, 0x13, 0x84, 0x69, 0x48, 0x7c, 0xaa, 0xff, 0x95, 0x40, 0xed, 0x87,
	0xf6, 0xdc, 0xa7, 0x41, 0x72, 0x49, 0x92, 0xdb, 0x41, 0x42, 0x7d, 0xf4, 0x3e, 0xb7, 0x6b, 0x7a,
	0xae, 0xe6, 0xc7, 0xe4, 0xec, 0xde, 0xed, 0x82, 0xf2, 0x40, 0xbc, 0xf9, 0xea, 0x3b, 0xe2, 0xc2,
	0xac, 0x6e, 0xe0, 0xd0, 0x05, 0xdf, 0xa5, 0x06, 0x16, 0x17, 0xfd, 0x47, 0xda, 0x2e, 0x80, 0xd2,
	0x85, 0x79, 0xd1, 0x35, 0xb1, 0x2a, 0xa1, 0x1d, 0xd8, 0x12, 0xe7, 0x89, 0x71, 0x35, 0x36, 0xf1,
	0xe0, 0xea, 0x8b, 0x5a, 0x60, 0xcd, 0xe0, 0xf5, 0x4e, 0x06, 0xc3, 0xbe, 0xf9, 0x4d, 0x2d, 0x22,
	0x0d, 0x76, 0x33, 0x86, 0x0d, 0x55, 0x66, 0xfe, 0xfd, 0xd1, 0x75, 0xd7, 0x32, 0x37, 0x46, 0x45,
	0xff, 0x23, 0x81, 0xfa, 0x78, 0xdd, 0x51, 0x17, 0x1a, 0x4e, 0x5a, 0xc8, 0x24, 0x22, 0xc9, 0xad,
	0x26, 0x35, 0x8b, 0x4f, 0x44, 0xf2, 0xb8, 0x54, 0x5c, 0x77, 0x32, 0x96, 0x75, 0xbf, 0x45, 0xb1,
	0xa2, 0xdf, 0x47, 0x00, 0x09, 0xd3, 0x89, 0x68, 0xb7, 0x10, 0x4f, 0x95, 0x5b, 0x86, 0xcf, 0x8c,
	0x43, 0x7e, 0x32, 0x0e, 0x07, 0xea, 0x59, 0x21, 0xa0, 0x93, 0xcc, 0x4c, 0x6b, 0x9d, 0xfd, 0x5c,
	0x7a, 0x19, 0xf5, 0x8a, 0x8f, 0xbf, 0x01, 0x25, 0x22, 0x31, 0xf1, 0xb5, 0x02, 0x2f, 0x66, 0xfb,
	0xc9, 0x1b, 0x81, 0x05, 0xae, 0x7f, 0x82, 0xca, 0x4a, 0x33, 0xcf, 0x6e, 0xcd, 0x7f, 0x07, 0xfa,
	0x25, 0x41, 0x49, 0x08, 0x06, 0x9d, 0x42, 0x71, 0xea, 0x39, 0x69, 0x1f, 0x8f, 0x9f, 0x91, 0x54,
	0xfa, 0x77, 0xe6, 0x52, 0xcf, 0xc1, 0x8c, 0x7b, 0x78, 0x0e, 0xb5, 0x8c, 0x0d, 0xa9, 0x50, 0xfc,
	0x49, 0x97, 0x69, 0x22, 0xec, 0xc8, 0xf2, 0x58, 0xed, 0x53, 0xe1, 0x85, 0x3c, 0x38, 0xae, 0xbf,
	0x03, 0x85, 0x6b, 0x71, 0xe3, 0x21, 0xbd, 0x98, 0x39, 0xc7, 0xbb, 0x07, 0xb0, 0x6f, 0x87, 0x7e,
	0x9b, 0xbf, 0xba, 0x6d, 0xfb, 0xae, 0xbd, 0x10, 0xef, 0xed, 0xcd, 0x7c, 0xfa, 0x2f, 0x00, 0x00,
	0xff, 0xff, 0x7d, 0xa8, 0x80, 0xe7, 0xab, 0x05, 0x00, 0x00,
}
