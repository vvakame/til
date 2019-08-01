package metago

func ValueOf(interface{}) Value {
	panic("in meta context!")
}

type Value interface {
	Fields() []Field
}

type Field interface {
	Name() string
	StructTag() StructTag
	Value() interface{}

	MatchTypeOf(typeHint TypeHint) bool
}

type StructTag interface {
	Get(string) string
}

type TypeHint struct {
	Receiver interface{}
}
