package dev_to_feed_test

import (
	"fmt"
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/dev_to/feed"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(dev_to_feed.New, nil)
	as.Nil(err)

	as.Equal("DEV Community - Feed", feed.Title)
	as.Equal("https://dev.to", feed.Link)
	// spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
	fmt.Println(feed.Items[0].Description)
}
