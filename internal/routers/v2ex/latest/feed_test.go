package v2ex_latest_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/v2ex/latest"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(v2ex_latest.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("V2EX - 全部主题", feed.Title)
	as.Equal("https://www.v2ex.com/?tab=all", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
