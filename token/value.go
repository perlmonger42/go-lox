package token

import "fmt"

type Value interface {
	Show() string   // converts strings to `"the string content"`
	String() string // converts strings to `the string content`
	TypeName() string

	IsNumber() bool
	IsString() bool
	IsBoolean() bool
	IsNil() bool
	IsObject() bool

	IsEqualTo(Value) bool
	equalsString(StringValue) bool
	equalsNumber(NumberValue) bool
	equalsBoolean(BooleanValue) bool
	equalsNil(NilValue) bool
	equalsObject(ObjectValue) bool
}

type StringValue struct {
	V string
}

var _ Value = StringValue{}

func (x StringValue) Show() string   { return fmt.Sprintf("%q", x.V) }
func (x StringValue) String() string { return fmt.Sprintf("%s", x.V) }

func (x StringValue) TypeName() string                  { return "string" }
func (x StringValue) IsNumber() bool                    { return false }
func (x StringValue) IsString() bool                    { return true }
func (x StringValue) IsBoolean() bool                   { return false }
func (x StringValue) IsNil() bool                       { return false }
func (x StringValue) IsObject() bool                    { return false }
func (x StringValue) IsEqualTo(v Value) bool            { return v.equalsString(x) }
func (x StringValue) equalsString(v StringValue) bool   { return v.V == x.V }
func (x StringValue) equalsNumber(v NumberValue) bool   { return false }
func (x StringValue) equalsBoolean(v BooleanValue) bool { return false }
func (x StringValue) equalsNil(v NilValue) bool         { return false }
func (x StringValue) equalsObject(v ObjectValue) bool   { return false }

type NumberValue struct {
	V float64
}

var _ Value = NumberValue{}

func (x NumberValue) Show() string   { return fmt.Sprintf("%g", x.V) }
func (x NumberValue) String() string { return fmt.Sprintf("%g", x.V) }

func (x NumberValue) TypeName() string                  { return "number" }
func (x NumberValue) IsNumber() bool                    { return true }
func (x NumberValue) IsString() bool                    { return false }
func (x NumberValue) IsBoolean() bool                   { return false }
func (x NumberValue) IsNil() bool                       { return false }
func (x NumberValue) IsObject() bool                    { return false }
func (x NumberValue) IsEqualTo(v Value) bool            { return v.equalsNumber(x) }
func (x NumberValue) equalsString(v StringValue) bool   { return false }
func (x NumberValue) equalsNumber(v NumberValue) bool   { return v.V == x.V }
func (x NumberValue) equalsBoolean(v BooleanValue) bool { return false }
func (x NumberValue) equalsNil(v NilValue) bool         { return false }
func (x NumberValue) equalsObject(v ObjectValue) bool   { return false }

type BooleanValue struct {
	V bool
}

var _ Value = BooleanValue{}

func (x BooleanValue) Show() string {
	if x.V {
		return "true"
	}
	return "false"
}
func (x BooleanValue) String() string {
	if x.V {
		return "true"
	}
	return "false"
}

func (x BooleanValue) TypeName() string                  { return "boolean" }
func (x BooleanValue) IsNumber() bool                    { return false }
func (x BooleanValue) IsString() bool                    { return false }
func (x BooleanValue) IsBoolean() bool                   { return true }
func (x BooleanValue) IsNil() bool                       { return false }
func (x BooleanValue) IsObject() bool                    { return false }
func (x BooleanValue) IsEqualTo(v Value) bool            { return v.equalsBoolean(x) }
func (x BooleanValue) equalsString(v StringValue) bool   { return false }
func (x BooleanValue) equalsNumber(v NumberValue) bool   { return false }
func (x BooleanValue) equalsBoolean(v BooleanValue) bool { return v.V == x.V }
func (x BooleanValue) equalsNil(v NilValue) bool         { return false }
func (x BooleanValue) equalsObject(v ObjectValue) bool   { return false }

type NilValue struct {
}

var _ Value = NilValue{}

func (x NilValue) Show() string   { return "nil" }
func (x NilValue) String() string { return "nil" }

func (x NilValue) TypeName() string                  { return "nil" }
func (x NilValue) IsNumber() bool                    { return false }
func (x NilValue) IsString() bool                    { return false }
func (x NilValue) IsBoolean() bool                   { return false }
func (x NilValue) IsNil() bool                       { return true }
func (x NilValue) IsObject() bool                    { return false }
func (x NilValue) IsEqualTo(v Value) bool            { return v.equalsNil(x) }
func (x NilValue) equalsString(v StringValue) bool   { return false }
func (x NilValue) equalsNumber(v NumberValue) bool   { return false }
func (x NilValue) equalsBoolean(v BooleanValue) bool { return false }
func (x NilValue) equalsNil(v NilValue) bool         { return true }
func (x NilValue) equalsObject(v ObjectValue) bool   { return true }

type Object interface {
	Show() string
	String() string
	EqualsObject(o Object) bool
}

type ObjectValue struct {
	V Object
}

var _ Value = ObjectValue{}

func (x ObjectValue) Show() string   { return x.V.Show() }
func (x ObjectValue) String() string { return x.V.String() }

func (x ObjectValue) TypeName() string                  { return "object" }
func (x ObjectValue) IsNumber() bool                    { return false }
func (x ObjectValue) IsString() bool                    { return false }
func (x ObjectValue) IsBoolean() bool                   { return false }
func (x ObjectValue) IsNil() bool                       { return false }
func (x ObjectValue) IsObject() bool                    { return true }
func (x ObjectValue) IsEqualTo(v Value) bool            { return v.equalsObject(x) }
func (x ObjectValue) equalsString(v StringValue) bool   { return false }
func (x ObjectValue) equalsNumber(v NumberValue) bool   { return false }
func (x ObjectValue) equalsBoolean(v BooleanValue) bool { return false }
func (x ObjectValue) equalsNil(v NilValue) bool         { return false }
func (x ObjectValue) equalsObject(v ObjectValue) bool   { return x.V == v.V }
