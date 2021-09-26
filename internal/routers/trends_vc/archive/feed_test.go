package trends_vc_archive_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	trends_vc_archive "github.com/chyroc/grss/internal/routers/trends_vc/archive"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(trends_vc_archive.New, nil)
	as.Nil(err)

	as.Equal("Trends.vc", feed.Title)
	as.Equal("Markets and Ideas", feed.Description)
	as.Equal("https://trends.vc/archives/", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
