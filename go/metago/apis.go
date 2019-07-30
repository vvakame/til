package metago

func ValueOf(interface{}) Value {
	panic("in meta context!")
}

func Types(...TypeArgument) GenericType {
	panic("in meta context!")
}

type Value interface {
	Fields() []Field
}

type Field interface {
	Name() string
	StructTag() StructTag
	Value() interface{}

	MatchTypeOf(typeHint TypeArgument) bool
}

type StructTag interface {
	Get(string) string
}

type GenericType int

type TypeArgument interface {
	isTypeArgument()
}

type TypeHint struct {
	Receiver interface{}
}

func (s TypeHint) isTypeArgument() {
	panic("in meta context!")
}
