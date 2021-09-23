package sspai_column_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	sspai_column "github.com/chyroc/grss/internal/routers/sspai/column"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(sspai_column.New, map[string]string{"id": "266"})
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("少数派专栏 - 生产力周报", feed.Title)
	as.Equal("https://sspai.com/column/266", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
