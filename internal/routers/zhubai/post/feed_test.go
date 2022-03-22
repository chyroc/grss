package zhubai_post_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	zhubai_post "github.com/chyroc/grss/internal/routers/zhubai/post"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(zhubai_post.New, map[string]string{"r": "via"})
	as.Nil(err)
	as.NotNil(feed)

	as.Equal("竹白 - 事不过三", feed.Title)
	as.Equal("https://via.zhubai.love/", feed.Link)
	as.Contains(feed.Description, "重要的事情不过三件：认识自己、好好学习、好好生活")
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
