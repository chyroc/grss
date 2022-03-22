package draveness_index_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	draveness_index "github.com/chyroc/grss/internal/routers/draveness/index"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(draveness_index.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	spew.Dump(feed, err)

	as.Equal("面向信仰编程 - 首页", feed.Title)
	as.Equal("https://draveness.me/index", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
