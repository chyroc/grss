package routers

import (
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/banyuetan/jinritan"
	"github.com/chyroc/grss/internal/routers/dev_to/feed"
	"github.com/chyroc/grss/internal/routers/lobste/home"
	"github.com/chyroc/grss/internal/routers/pingwest/status"
	"github.com/chyroc/grss/internal/routers/reddit/community/hot"
	"github.com/chyroc/grss/internal/routers/sspai/column"
	"github.com/chyroc/grss/internal/routers/sspai/matrix"
	"github.com/chyroc/grss/internal/routers/trends_vc/archive"
	"github.com/chyroc/grss/internal/routers/v2ex/latest"
	"github.com/chyroc/grss/internal/routers/weibo/user/origin"
	"github.com/chyroc/grss/internal/routers/woshipm/latest"
	"github.com/chyroc/grss/internal/routers/xueqiu/livenews"
	"github.com/chyroc/grss/internal/routers/xueqiu/snb_article"
	"github.com/chyroc/grss/internal/routers/ycombinator/best"
	"github.com/chyroc/grss/internal/routers/ycombinator/newest"
	"github.com/chyroc/grss/internal/routers/zhihu/bookstore/newest"
)

func init() {
	Register("/banyuetan/jinritan", banyuetan_jinritan.New)
	Register("/dev_to/feed", dev_to_feed.New)
	Register("/lobste/home", lobste_home.New)
	Register("/pingwest/status", pingwest_status.New)
	Register("/reddit/community/hot/:r", reddit_community_hot.New,
		map[string]string{"r": "golang"},
	)
	Register("/sspai/column/:id", sspai_column.New,
		map[string]string{"id": "264"},
		map[string]string{"id": "266"},
	)
	Register("/sspai/matrix", sspai_matrix.New)
	Register("/trends_vc/archive", trends_vc_archive.New)
	Register("/v2ex/latest", v2ex_latest.New)
	Register("/weibo/user/origin/:uid", weibo_user_origin.New,
		map[string]string{"uid": "1088413295"},
	)
	Register("/woshipm/latest", woshipm_latest.New)
	Register("/xueqiu/livenews", xueqiu_livenews.New)
	Register("/xueqiu/snb_article", xueqiu_snb_article.New)
	Register("/ycombinator/best", ycombinator_best.New)
	Register("/ycombinator/newest", ycombinator_newest.New)
	Register("/zhihu/bookstore/newest", zhihu_bookstore_newest.New)
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
