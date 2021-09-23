package objx_test

import (
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/helper/objx"
)

func Test_ToList(t *testing.T) {
	as := assert.New(t)

	resp := []*item{}
	err := objx.
		New([]string{"1", "2"}).
		Map(func(idx int, v interface{}) interface{} {
			return &item{Name: v.(string)}
		}).
		ToList(&resp)
	as.Nil(err)
	as.Len(resp, 2)
	as.Equal("1", resp[0].Name)
	as.Equal("2", resp[1].Name)
}
