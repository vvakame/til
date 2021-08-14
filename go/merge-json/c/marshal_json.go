package c

import (
	"encoding/json"
	"errors"
	"reflect"
)

func MarshalFlatten(objs ...interface{}) ([]byte, error) {
	if len(objs) == 0 {
		return []byte("{}"), nil
	} else if len(objs) == 1 {
		return json.Marshal(objs[0])
	}

	merged, err := mergeFields(objs...)
	if err != nil {
		return nil, err
	}

	return json.Marshal(merged)
}

func mergeFields(objs ...interface{}) (interface{}, error) {
	var sfs []reflect.StructField
	var vs []reflect.Value
	for _, obj := range objs {
		v := reflect.ValueOf(obj)
		if v.Kind() != reflect.Ptr {
			return nil, errors.New("must be a pointer")
		}
		if v.IsNil() {
			return nil, errors.New("must be a non-nil pointer")
		}
		v = v.Elem()
		if v.Kind() != reflect.Struct {
			return nil, errors.New("must be a pointer to struct")
		}

		for i := 0; i < v.NumField(); i++ {
			sf := v.Type().Field(i)
			v := v.Field(i)
			if !v.CanSet() {
				continue
			}

			sfs = append(sfs, sf)
			vs = append(vs, v)
		}
	}

	st := reflect.StructOf(sfs)

	s := reflect.New(st).Elem()
	for i := 0; i < len(vs); i++ {
		s.Field(i).Set(vs[i])
	}

	return s.Interface(), nil
}
