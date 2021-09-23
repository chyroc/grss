package objx_test

import (
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/helper/objx"
)

type item struct {
	Name string
}

func Test_Map(t *testing.T) {
	as := assert.New(t)

	var data interface{} = []*item{
		{Name: "a"},
		{Name: "b"},
	}

	obj := objx.New(data).Map(func(idx int, v interface{}) interface{} {
		return v.(*item).Name
	}).Obj()
	as.Equal([]interface{}{"a", "b"}, obj)
}
