package routers

import (
	"github.com/chyroc/grss/internal/fetch"
	banyuetan_jinritan "github.com/chyroc/grss/internal/routers/banyuetan/jinritan"
	"github.com/chyroc/grss/internal/routers/pingwest/status"
	"github.com/chyroc/grss/internal/routers/sspai/column"
	"github.com/chyroc/grss/internal/routers/sspai/matrix"
	"github.com/chyroc/grss/internal/routers/zhihu/bookstore/zhihu_bookstore_newest"
)

func init() {
	Register("/zhihu/bookstore/newest", zhihu_bookstore_newest.New)
	Register("/pingwest/status", pingwest_status.New)
	Register("/sspai/matrix", sspai_matrix.New)
	Register("/sspai/column/:id", sspai_column.New,
		map[string]string{"id": "264"},
		map[string]string{"id": "266"},
	)
	Register("/banyuetan/jinritan", banyuetan_jinritan.New)
}

var routers []Router

type Router struct {
	Path   string
	Source func(map[string]string) (*fetch.Source, error)
	Args   []map[string]string
}

func Register(path string, source func(map[string]string) (*fetch.Source, error), args ...map[string]string) {
	if len(args) == 0 {
		args = []map[string]string{nil}
	}
	routers = append(routers, Router{
		Path:   path,
		Source: source,
		Args:   args,
	})
}

func Get() []Router {
	return routers
}
