package jj

// ElemType element type
type ElemType int

const (
	// StringType
	StringType ElemType = iota + 1
	// NumberType
	NumberType
	// BoolType
	BoolType
	// NullType
	NullType
	// ArrayType
	ArrayType
	// ObjectType
	ObjectType
)

// String
func (e ElemType) String() string {
	switch e {
	case StringType:
		return "string"
	case NumberType:
		return "number"
	case BoolType:
		return "bool"
	case NullType:
		return "null"
	case ArrayType:
		return "array"
	case ObjectType:
		return "object"
	default:
		return "unknown"
	}
}

type (
	// Elem
	Elem interface {
		Type() ElemType
		Val() any
		Get(path ...any) Elem
		AsString() string
		AsNumber() float64
		AsBool() bool
		AsNull()
		AsArray() []Elem
		AsObject() map[string]Elem
	}

	// StringElem
	StringElem struct {
		val string
	}

	// NumberElem
	NumberElem struct {
		val float64
	}

	BoolElem struct {
		val bool
	}

	NullElem struct {
	}

	ArrayElem struct {
		val []Elem
	}

	ObjectElem struct {
		val map[string]Elem
	}
)

func NewLiteralElem(val any) Elem {
	if val == nil {
		return &NullElem{}
	}
	switch v := val.(type) {
	case string:
		return &StringElem{val: v}
	case float64:
		return &NumberElem{val: v}
	case bool:
		return &BoolElem{val: v}
	default:
		panic("not support")
	}
}

func (e *StringElem) Get(...any) Elem {
	return e
}

func (e *NumberElem) Get(...any) Elem {
	return e
}

func (e *BoolElem) Get(...any) Elem {
	return e
}

func (e *NullElem) Get(...any) Elem {
	return e
}

func (e *ArrayElem) Get(path ...any) Elem {
	i := path[0].(int)
	return e.val[i].Get(path[1:]...)
}

func (e *ObjectElem) Get(path ...any) Elem {
	key := path[0].(string)
	return e.val[key].Get(path[1:]...)
}

func (e *ObjectElem) AsString() string {
	return ""
}

func (e *ObjectElem) AsNumber() float64 {
	return 0
}

func (e *ObjectElem) AsBool() bool {
	return false
}

func (e *ObjectElem) AsNull() {
}

func (e *ObjectElem) AsArray() []Elem {
	return nil
}

func (e *ObjectElem) AsObject() map[string]Elem {
	return e.val
}

func (e *ArrayElem) AsString() string {
	return ""
}

func (e *ArrayElem) AsNumber() float64 {
	return 0
}

func (e *ArrayElem) AsBool() bool {
	return false
}

func (e *ArrayElem) AsNull() {
}

func (e *ArrayElem) AsArray() []Elem {
	return e.val
}

func (e *ArrayElem) AsObject() map[string]Elem {
	return nil
}

func (e *NullElem) AsString() string {
	return ""
}

func (e *NullElem) AsNumber() float64 {
	return 0
}

func (e *NullElem) AsBool() bool {
	return false
}

func (e *NullElem) AsNull() {
}

func (e *NullElem) AsArray() []Elem {
	return nil
}

func (e *NullElem) AsObject() map[string]Elem {
	return nil
}

func (e *BoolElem) AsString() string {
	return ""
}

func (e *BoolElem) AsNumber() float64 {
	return 0
}

func (e *BoolElem) AsBool() bool {
	return e.val
}

func (e *BoolElem) AsNull() {
}

func (e *BoolElem) AsArray() []Elem {
	return nil
}

func (e *BoolElem) AsObject() map[string]Elem {
	return nil
}

func (e *NumberElem) AsString() string {
	return ""
}

func (e *NumberElem) AsNumber() float64 {
	return e.val
}

func (e *NumberElem) AsBool() bool {
	return false
}

func (e *NumberElem) AsNull() {
}

func (e *NumberElem) AsArray() []Elem {
	return nil
}

func (e *NumberElem) AsObject() map[string]Elem {
	return nil
}

func (e *StringElem) AsString() string {
	return e.val
}

func (e *StringElem) AsNumber() float64 {
	return 0
}

func (e *StringElem) AsBool() bool {
	return false
}

func (e *StringElem) AsNull() {
}

func (e *StringElem) AsArray() []Elem {
	return nil
}

func (e *StringElem) AsObject() map[string]Elem {
	return nil
}

func NewArrayElem(val []Elem) *ArrayElem {
	return &ArrayElem{val: val}
}

func NewObjectElem(val map[string]Elem) *ObjectElem {
	return &ObjectElem{val: val}
}

func (e *StringElem) Type() ElemType {
	return StringType
}

func (e *NumberElem) Type() ElemType {
	return NumberType
}

func (e *BoolElem) Type() ElemType {
	return BoolType
}

func (e *ArrayElem) Type() ElemType {
	return ArrayType
}

func (e *ObjectElem) Type() ElemType {
	return ObjectType
}

func (e *NullElem) Type() ElemType {
	return NullType
}

func (e *StringElem) Val() any {
	return e.val
}

func (e *NumberElem) Val() any {
	return e.val
}

func (e *BoolElem) Val() any {
	return e.val
}

func (e *ArrayElem) Val() any {
	return e.val
}

func (e *ObjectElem) Val() any {
	return e.val
}

func (e *NullElem) Val() any {
	return nil
}
