package sspai_matrix_test

import (
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/fetch"
	"github.com/chyroc/grss/interface/routers/sspai/matrix"
	"github.com/davecgh/go-spew/spew"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(sspai_matrix.New())
	as.Nil(err)

	as.Equal("少数派 - Matrix", feed.Title)
	as.Equal("https://sspai.com/matrix", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
