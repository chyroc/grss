package routers

import (
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/banyuetan/jinritan"
	"github.com/chyroc/grss/internal/routers/dev_to/feed"
	"github.com/chyroc/grss/internal/routers/draveness/index"
	"github.com/chyroc/grss/internal/routers/github/trending"
	"github.com/chyroc/grss/internal/routers/gocn/daily"
	"github.com/chyroc/grss/internal/routers/lobste/home"
	"github.com/chyroc/grss/internal/routers/meituan_tech/article"
	"github.com/chyroc/grss/internal/routers/pingwest/status"
	"github.com/chyroc/grss/internal/routers/reddit/community/hot"
	"github.com/chyroc/grss/internal/routers/sspai/column"
	"github.com/chyroc/grss/internal/routers/sspai/index"
	"github.com/chyroc/grss/internal/routers/sspai/matrix"
	"github.com/chyroc/grss/internal/routers/studygolang/index"
	"github.com/chyroc/grss/internal/routers/todtod/index"
	"github.com/chyroc/grss/internal/routers/toutiaoio/index"
	"github.com/chyroc/grss/internal/routers/trends_vc/archive"
	"github.com/chyroc/grss/internal/routers/twitter/user/origin"
	"github.com/chyroc/grss/internal/routers/v2ex/latest"
	"github.com/chyroc/grss/internal/routers/weibo/user/origin"
	"github.com/chyroc/grss/internal/routers/woshipm/latest"
	"github.com/chyroc/grss/internal/routers/xueqiu/livenews"
	"github.com/chyroc/grss/internal/routers/xueqiu/snb_article"
	"github.com/chyroc/grss/internal/routers/ycombinator/best"
	"github.com/chyroc/grss/internal/routers/ycombinator/newest"
	"github.com/chyroc/grss/internal/routers/zhihu/bookstore/newest"
	"github.com/chyroc/grss/internal/routers/zhubai/post"
)

func init() {
	Register("/banyuetan/jinritan", banyuetan_jinritan.New)
	Register("/dev_to/feed", dev_to_feed.New)
	Register("/draveness/index", draveness_index.New)
	Register("/github/trending/:lang/:since", github_trending.New,
		map[string]string{"lang": "go", "since": "daily"},
		map[string]string{"lang": "python", "since": "daily"},
		map[string]string{"lang": "rust", "since": "daily"},
	)
	Register("/gocn/daily", gocn_daily.New)
	Register("/lobste/home", lobste_home.New)
	Register("/meituan_tech/article", meituan_tech_article.New)
	Register("/pingwest/status", pingwest_status.New)
	Register("/reddit/community/hot/:r", reddit_community_hot.New,
		map[string]string{"r": "golang"},
		map[string]string{"r": "tech"},
		map[string]string{"r": "technology"},
		map[string]string{"r": "geek"},
	)
	Register("/sspai/column/:id", sspai_column.New,
		map[string]string{"id": "264"},
		map[string]string{"id": "266"},
	)
	Register("/sspai/index", sspai_index.New)
	Register("/sspai/matrix", sspai_matrix.New)
	Register("/studygolang/index", studygolang_index.New)
	Register("/todtod/index", todtod_index.New)
	Register("/toutiaoio/index", toutiaoio_index.New)
	Register("/trends_vc/archive", trends_vc_archive.New)
	Register("/twitter/user/origin/:uid", twitter_user_origin.New,
		map[string]string{"uid": "woshisuchang"},
	)
	Register("/v2ex/latest", v2ex_latest.New)
	Register("/weibo/user/origin/:uid", weibo_user_origin.New,
		map[string]string{"uid": "1088413295"},
		map[string]string{"uid": "5722964389"},
	)
	Register("/woshipm/latest", woshipm_latest.New)
	Register("/xueqiu/livenews", xueqiu_livenews.New)
	Register("/xueqiu/snb_article", xueqiu_snb_article.New)
	Register("/ycombinator/best", ycombinator_best.New)
	Register("/ycombinator/newest", ycombinator_newest.New)
	Register("/zhihu/bookstore/newest", zhihu_bookstore_newest.New)
	Register("/zhubai/post/:r", zhubai_post.New,
		map[string]string{"r": "chasays"},
		map[string]string{"r": "pythonhunter"},
	)
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
