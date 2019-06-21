//+build !metago

package govalidtemplate

import (
	"bytes"
	"encoding/json"
)

func (d *Data) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")
	{
		buf.WriteString(`"`)
		buf.WriteString("ID")
		buf.WriteString(`":`)

		b, err := json.Marshal(d.ID)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.WriteString(",")
	{
		buf.WriteString(`"`)
		buf.WriteString("Name")
		buf.WriteString(`":`)

		b, err := json.Marshal(d.Name)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.WriteString(",")
	{
		buf.WriteString(`"`)
		buf.WriteString("Age")
		buf.WriteString(`":`)

		b, err := json.Marshal(d.Age)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}
