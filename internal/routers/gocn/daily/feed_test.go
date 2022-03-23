package gocn_daily_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	gocn_daily "github.com/chyroc/grss/internal/routers/gocn/daily"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(gocn_daily.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	spew.Dump(feed, err)

	as.Equal("GoCN - 每日新闻", feed.Title)
	as.Equal("https://gocn.vip/topics/cate/18", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
