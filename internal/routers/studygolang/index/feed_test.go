package studygolang_index_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	studygolang_index "github.com/chyroc/grss/internal/routers/studygolang/index"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(studygolang_index.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	spew.Dump(feed, err)

	as.Equal("Go语言中文网 - 首页", feed.Title)
	as.Equal("https://studygolang.com/", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
