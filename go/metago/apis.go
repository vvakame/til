package metago

func ValueOf(interface{}) Value {
	panic("in meta context!")
}

type Value interface {
	Fields() []Field
}

type Field interface {
	Name() string
	StructTagGet(string) string
	Value() interface{}

	MatchTypeOf(typeHint TypeHint) bool
}

type TypeHint struct {
	Receiver interface{}
}
