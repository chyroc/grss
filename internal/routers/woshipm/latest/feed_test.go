package woshipm_latest_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	woshipm_latest "github.com/chyroc/grss/internal/routers/woshipm/latest"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(woshipm_latest.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("人人都是产品经理 - 最新文章", feed.Title)
	as.Equal("http://www.woshipm.com", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
