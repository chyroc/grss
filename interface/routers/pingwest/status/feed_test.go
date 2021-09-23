package pingwest_status_test

import (
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/fetch"
	"github.com/chyroc/grss/interface/routers/pingwest/status"
	"github.com/davecgh/go-spew/spew"
)

func TestName(t *testing.T) {
	as := assert.New(t)
	as.Nil(nil)

	feed, err := fetch.Fetch(pingwest_status.New())
	as.Nil(err)
	as.Equal("品玩 - 实时要闻", feed.Title)
	as.Equal("https://www.pingwest.com/status", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
}
