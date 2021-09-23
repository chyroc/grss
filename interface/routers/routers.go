package routers

import (
	"github.com/chyroc/grss/interface/fetch"
	"github.com/chyroc/grss/interface/routers/pingwest/status"
	"github.com/chyroc/grss/interface/routers/zhihu/bookstore/zhihu_bookstore_newest"
)

func init() {
	Register("/zhihu/bookstore/newest", zhihu_bookstore_newest.New())
	Register("/pingwest/status", pingwest_status.New())
}

var routers []Router

type Router struct {
	Path   string
	Source fetch.Source
}

func Register(path string, source fetch.Source) {
	routers = append(routers, Router{
		Path:   path,
		Source: source,
	})
}

func Get() []Router {
	return routers
}
