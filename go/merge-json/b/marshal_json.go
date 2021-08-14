package b

import (
	"encoding/json"
	"github.com/imdario/mergo"
)

func MarshalFlatten(objs ...interface{}) ([]byte, error) {
	if len(objs) == 0 {
		return []byte("{}"), nil
	} else if len(objs) == 1 {
		return json.Marshal(objs[0])
	}

	dst := make(map[string]interface{})

	for _, obj := range objs {
		err := mergo.Map(&dst, obj, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(dst)
}
