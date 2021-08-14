package c

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

var m = &merger{
	cacheValue: &cacheValue{
		nextTypes: make(map[reflect.Type]*cacheValue),
	},
}

func MarshalFlatten(objs ...interface{}) ([]byte, error) {
	if len(objs) == 0 {
		return []byte("{}"), nil
	} else if len(objs) == 1 {
		return json.Marshal(objs[0])
	}

	merged, err := m.mergeObjects(objs...)
	if err != nil {
		return nil, err
	}

	return json.Marshal(merged)
}

type merger struct {
	*cacheValue
}

type cacheValue struct {
	sync.RWMutex
	nextTypes map[reflect.Type]*cacheValue
	typeCache *typeCache
}

type typeCache struct {
	st  reflect.Type
	vas []*valueAccessor
}

type valueAccessor struct {
	skip bool

	toFieldIndex int

	objIndex   int
	fieldIndex int
}

func (m *merger) mergeObjects(objs ...interface{}) (interface{}, error) {
	typeCache, err := m.getTypeCache(objs, nil)
	if err != nil {
		return nil, err
	}

	return typeCache.mergeObjects(objs...)
}

func (cv *cacheValue) getTypeCache(objs, rest []interface{}) (*typeCache, error) {
	if len(objs) == 0 {
		return nil, errors.New("objs len is 0")
	}
	if len(rest) == 0 {
		rest = objs
	}

	obj := rest[0]
	v, err := toBareStructValue(obj)
	if err != nil {
		return nil, err
	}

	if len(rest) == 1 {
		cv.RLock()
		if cv.typeCache != nil {
			cv.RUnlock()
			return cv.typeCache, nil
		}
		cv.RUnlock()

		cv.Lock()

		var sfs []reflect.StructField
		var vas []*valueAccessor
		var currentLoop int
		for objIdx, obj := range objs {
			v, err := toBareStructValue(obj)
			if err != nil {
				return nil, err
			}

			for i := 0; i < v.NumField(); i++ {
				sf := v.Type().Field(i)
				v := v.Field(i)

				if !v.CanSet() {
					vas = append(vas, &valueAccessor{
						skip: true,
					})
					continue
				}

				sfs = append(sfs, sf)
				vas = append(vas, &valueAccessor{
					toFieldIndex: currentLoop,
					objIndex:     objIdx,
					fieldIndex:   i,
				})
				currentLoop++
			}
		}

		st := reflect.StructOf(sfs)

		cv.typeCache = &typeCache{
			st:  st,
			vas: vas,
		}
		cv.Unlock()

		return cv.typeCache, nil
	}

	cv.RLock()
	v, err = toBareStructValue(rest[1])
	if err != nil {
		return nil, err
	}
	next, ok := cv.nextTypes[v.Type()]
	if !ok {
		cv.RUnlock()
		cv.Lock()
		next = &cacheValue{
			nextTypes: make(map[reflect.Type]*cacheValue),
		}
		cv.nextTypes[v.Type()] = next
		cv.Unlock()
	}

	return next.getTypeCache(objs, rest[1:])
}

func (tc *typeCache) mergeObjects(objs ...interface{}) (interface{}, error) {
	s := reflect.New(tc.st).Elem()
	for _, va := range tc.vas {
		if va.skip {
			continue
		}

		obj := objs[va.objIndex]

		v, err := toBareStructValue(obj)
		if err != nil {
			return nil, err
		}

		s.Field(va.toFieldIndex).Set(v.Field(va.fieldIndex))
	}

	return s.Interface(), nil
}

func toBareStructValue(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("must be a pointer")
	}
	if v.IsNil() {
		return reflect.Value{}, errors.New("must be a non-nil pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("must be a pointer to struct")
	}

	return v, nil
}
