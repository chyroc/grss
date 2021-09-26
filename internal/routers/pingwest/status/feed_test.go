package pingwest_status_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/pingwest/status"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	as := assert.New(t)
	as.Nil(nil)

	feed, err := fetch.Fetch(pingwest_status.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("品玩 - 实时要闻", feed.Title)
	as.Equal("https://www.pingwest.com/status", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
}
