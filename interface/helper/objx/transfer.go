package objx

func (r *Obj) Transfer(f func(obj interface{}) interface{}) *Obj {
	if r.err != nil {
		return r
	}

	r.obj = f(r.obj)

	return r
}
