package twitter_user_origin_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	twitter_user_origin "github.com/chyroc/grss/internal/routers/twitter/user/origin"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(twitter_user_origin.New, map[string]string{"uid": "awscloud"})
	as.Nil(err)

	as.Equal("Twitter - Amazon Web Services Origin Twitter", feed.Title)
	as.Equal("https://twitter.com/awscloud/", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	// as.NotEmpty(feed.Items[0].Description)
}
