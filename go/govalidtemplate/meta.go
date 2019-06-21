package govalidtemplate

func metaContext() *MetaContext {
	panic("🐾")
}

type MetaContext struct {
}

type MetaReceiver struct {
	Fields []*MetaField
}

type MetaField struct {
	Name string
}

type MetaRef struct {
}

func (mc *MetaContext) Get(v interface{}) *MetaReceiver {
	panic("🐾")
}

func (mr *MetaReceiver) Field(field *MetaField) *MetaRef {
	panic("🐾")
}

func (mr *MetaReceiver) IsLastField(field *MetaField) bool {
	panic("🐾")
}
