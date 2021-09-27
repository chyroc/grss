package xueqiu_livenews_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/xueqiu/livenews"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(xueqiu_livenews.New, nil)
	as.Nil(err)

	as.Equal("雪球 - 7X24快讯", feed.Title)
	as.Equal("https://xueqiu.com/?category=livenews", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
