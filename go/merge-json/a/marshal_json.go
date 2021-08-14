package a

import (
	"bytes"
	"encoding/json"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func MarshalFlatten(objs ...interface{}) ([]byte, error) {
	if len(objs) == 0 {
		return []byte("{}"), nil
	}

	// めんどくさいので要素がarrayの場合の考慮は一旦しない…

	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	if len(objs) == 1 {
		err := json.NewEncoder(buf).Encode(objs[0])
		return buf.Bytes(), err
	}

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