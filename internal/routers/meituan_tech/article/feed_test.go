package meituan_tech_article_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/meituan_tech/article"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch(meituan_tech_article.New, nil)
	as.Nil(err)
	as.NotNil(feed)

	spew.Dump(feed, err)

	as.Equal("美团技术团队 - 文章", feed.Title)
	as.Equal("https://tech.meituan.com/", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
