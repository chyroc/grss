package zhihu_bookstore_newest_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/zhihu/bookstore/newest"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_ZhihuBookstore(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(zhihu_bookstore_newest.New, nil)
	as.Nil(err)

	as.Equal("知乎书店 - 新书抢鲜", feed.Title)
	as.Equal("https://www.zhihu.com/pub/features/new", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
}
