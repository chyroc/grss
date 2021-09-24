package ycombinator_newest_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/ycombinator/newest"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(ycombinator_newest.New, nil)
	as.Nil(err)

	as.Equal("Hacker News - New Links", feed.Title)
	as.Equal("https://news.ycombinator.com/newest", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
