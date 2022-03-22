package todtod_index_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	todtod_index "github.com/chyroc/grss/internal/routers/todtod/index"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(todtod_index.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	spew.Dump(feed, err)

	as.Equal("TO-D 杂志 - 首页", feed.Title)
	as.Equal("https://todtod.io/", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
