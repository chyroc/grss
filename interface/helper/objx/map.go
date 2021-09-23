package objx

import (
	"fmt"
	"reflect"
)

func (r *Obj) Map(f func(idx int, v interface{}) interface{}) *Obj {
	if r.err != nil {
		return r
	}

	arr, err := interfaceToInterfaceList(r.obj)
	if err != nil {
		r.err = err
		return r
	}

	objs := []interface{}{}
	for i, v := range arr {
		res := f(i, v)
		objs = append(objs, res)
	}
	r.obj = objs
	return r
}

func interfaceToInterfaceList(v interface{}) (res []interface{}, err error) {
	vv := reflect.ValueOf(v)
	canToArrKinds := []reflect.Kind{
		reflect.String,
		reflect.Slice,
		reflect.Array,
	}
	if !isInKind(vv.Kind(), canToArrKinds) {
		return nil, fmt.Errorf("%T unsupport .Map", v)
	}
	switch vv.Kind() {
	case reflect.String:
		for _, v := range []rune(vv.String()) {
			res = append(res, rune(v))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < vv.Len(); i++ {
			res = append(res, vv.Index(i).Interface())
		}
	}
	return
}

func isInKind(kind reflect.Kind, kinds []reflect.Kind) bool {
	for _, v := range kinds {
		if kind == v {
			return true
		}
	}
	return false
}
