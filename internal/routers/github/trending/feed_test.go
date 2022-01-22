package github_trending_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	github_trending "github.com/chyroc/grss/internal/routers/github/trending"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	as := assert.New(t)
	as.Nil(nil)

	feed, err := fetch.Fetch(github_trending.New, map[string]string{
		"lang":  "go",
		"since": "daily",
	})
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("GitHub - Trending - go - daily", feed.Title)
	as.Equal("https://github.com", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
}
