package objx

import (
	"fmt"
	"reflect"
	"unsafe"
)

func (r *Obj) Obj() interface{} {
	return r.obj
}

func (r *Obj) StringList() (res []string, err error) {
	if r.err != nil {
		return nil, r.err
	}
	arr, err := interfaceToInterfaceList(r.obj)
	if err != nil {
		return nil, err
	}
	for _, v := range arr {
		res = append(res, fmt.Sprintf("%s", v))
	}
	return res, nil
}

func (r *Obj) ToList(resp interface{}) (err error) {
	if r.err != nil {
		return r.err
	}
	arr, err := interfaceToInterfaceList(r.obj)
	if err != nil {
		return err
	}
	respV := reflect.ValueOf(resp)
	respT := reflect.TypeOf(resp)
	if respV.Kind() != reflect.Ptr {
		return fmt.Errorf("resp must be ptr")
	}
	respV = respV.Elem()
	respT = respT.Elem()
	if respV.Kind() != reflect.Slice {
		return fmt.Errorf("resp must be ptr of slice")
	}
	for i := 0; i < len(arr); i++ {
		objV := reflect.NewAt(respT.Elem().Elem(), unsafe.Pointer(reflect.ValueOf(arr[i]).Elem().UnsafeAddr()))
		respV.Set(reflect.Append(respV, objV))
	}
	return nil
}
