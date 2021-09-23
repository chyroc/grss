package objx

func (r *Obj) Filter(f func(idx int, obj interface{}) bool) *Obj {
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
		if f(i, v) {
			objs = append(objs, v)
		}
	}
	r.obj = objs

	return r
}
