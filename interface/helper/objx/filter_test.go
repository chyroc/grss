package objx_test

import (
	"strconv"
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/helper/objx"
)

func Test_Filter(t *testing.T) {
	as := assert.New(t)

	res := objx.New([]string{"1", "2", "3"}).Filter(func(idx int, obj interface{}) bool {
		i, _ := strconv.ParseInt(obj.(string), 10, 64)
		return i%2 == 0
	}).Obj()
	as.Equal([]interface{}{"2"}, res)
}
