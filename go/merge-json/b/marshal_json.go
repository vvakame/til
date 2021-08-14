package b

import (
	"bytes"
	"encoding/json"
	"github.com/imdario/mergo"
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

	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	if len(objs) == 1 {
		err := json.NewEncoder(buf).Encode(objs[0])
		return buf.Bytes(), err
	}

	dst := make(map[string]interface{})

	for _, obj := range objs {
		err := mergo.Map(&dst, obj, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	err := json.NewEncoder(buf).Encode(dst)
	return buf.Bytes(), err
}
