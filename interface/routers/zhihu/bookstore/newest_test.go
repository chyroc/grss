package zhihu_bookstore_test

import (
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/fetch"
	zhihu_bookstore "github.com/chyroc/grss/interface/routers/zhihu/bookstore"
	"github.com/davecgh/go-spew/spew"
)

func Test_ZhihuBookstore(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(zhihu_bookstore.New())
	as.Nil(err)

	as.Equal("知乎书店-新书抢鲜", feed.Title)
	as.Equal("https://www.zhihu.com/pub/features/new", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
}
