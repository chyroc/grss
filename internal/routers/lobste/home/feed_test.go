package lobste_home_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	lobste_home "github.com/chyroc/grss/internal/routers/lobste/home"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	as := assert.New(t)
	as.Nil(nil)

	feed, err := fetch.Fetch(lobste_home.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("Lobsters - Home", feed.Title)
	as.Equal("https://lobste.rs/", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
}
