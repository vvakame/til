package a

import (
	"bytes"
	"encoding/json"
)

func MarshalFlatten(objs ...interface{}) ([]byte, error) {
	if len(objs) == 0 {
		return []byte("{}"), nil
	} else if len(objs) == 1 {
		return json.Marshal(objs[0])
	}

	// めんどくさいので要素がarrayの場合の考慮は一旦しない…

	var buf bytes.Buffer
	for idx, obj := range objs {
		b, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		if idx == 0 {
			// head
			buf.Write(b[0 : len(b)-1])
			buf.WriteString(",")
		} else if idx != len(objs)-1 {
			// body
			buf.Write(b[1 : len(b)-1])
			buf.WriteString(",")
		} else {
			// tail
			buf.Write(b[1:])
		}
	}

	return buf.Bytes(), nil
}
