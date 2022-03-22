package post_test

import (
	"net/http"
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	reddit_community_hot "github.com/chyroc/grss/internal/routers/reddit/community/hot"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(reddit_community_hot.New, map[string]string{"r": "golang"})
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("Reddit - The Go Programming Language", feed.Title)
	as.Equal("https://www.reddit.com/r/golang/hot/", feed.Link)
	as.Contains(feed.Description, "about the Go programming language and related tools, events")
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}

func Test_fetchcontent(t *testing.T) {
	as := assert.New(t)
	text, err := helper.Req.New(http.MethodGet, "https://www.reddit.com/r/golang/comments/pv7zez/go_proverbs/").Text()
	as.Nil(err)
	as.Contains(text, "Don't communicate by sharing memory, share memory by communicating.")
}
