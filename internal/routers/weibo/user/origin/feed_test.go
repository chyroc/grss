package weibo_user_origin_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/weibo/user/origin"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(weibo_user_origin.New, map[string]string{"uid": "5722964389"})
	as.Nil(err)

	as.Equal("微博 - Easy 原创微博", feed.Title)
	as.Equal("https://m.weibo.cn/p/1088413295", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
