//+build metago

package govalidtemplate

import (
	"bytes"
	"encoding/json"
)

func (d *Data) MarshalJSON() ([]byte, error) {
	// NOTE _ prefixの変数と、それが絡むstatementは全部消えたり置き換えられたりする

	_meta := metaContext()

	var buf bytes.Buffer
	buf.WriteString("{")

	_obj := _meta.Get(d)
	for _, _field := range _obj.Fields {
		buf.WriteString(`"`)
		buf.WriteString(_field.Name)
		buf.WriteString(`":`)

		b, err := json.Marshal(_obj.Field(_field))
		if err != nil {
			return nil, err
		}
		buf.Write(b)

		if !_obj.IsLastField(_field) {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}
