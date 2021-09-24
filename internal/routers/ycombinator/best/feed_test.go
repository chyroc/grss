package ycombinator_best_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/ycombinator/best"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(ycombinator_best.New, nil)
	as.Nil(err)

	as.Equal("Hacker News - Top Links", feed.Title)
	as.Equal("https://news.ycombinator.com/best", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
