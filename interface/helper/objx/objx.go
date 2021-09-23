package objx

type Obj struct {
	err     error
	obj     interface{}
	objType objType
}

func New(obj interface{}) *Obj {
	return &Obj{obj: obj, objType: objTypeObj}
}

type objType int

const (
	objTypeObj objType = iota
	objTypeArr
)
