package main

import (
	"os"
	"strings"

	"github.com/chyroc/grss/internal/helper"
)

func main() {
	url := os.Args[1]
	path := os.Args[2]
	pkgName := strings.ReplaceAll(path, "/", "_")

	assert(os.MkdirAll("internal/routers/"+path, 0o777))

	text, err := genFeed(url, path, pkgName)
	assert(err)
	assert(os.WriteFile("internal/routers/"+path+"/feed.go", []byte(text), 0o666))

	textTest, err := genTest(path, pkgName)
	assert(err)
	assert(os.WriteFile("internal/routers/"+path+"/feed_test.go", []byte(textTest), 0o666))
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func genFeed(url, path, pkgName string) (string, error) {
	return helper.BuildTemplate(`package {{ .PkgPath }}

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "{{ .Link }}"
	return &fetch.Source{
		Title:       "",
		Description: "",
		Link:        link,

		Fetch: func() (interface{}, error) {
			text, err := new(sspaiMatrixResp)
			text,err := helper.Req.New(http.MethodGet, link).Text()
			return text,err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(string)))
				if err != nil {
					return nil, err
				}
				itemSelections := helper.Selection2List(q.Find(`+"`"+`section.item`+"`"+`))

			err = lambda.New(itemSelections).MapArrayAsync(func(idx int, obj interface{}) interface{} {
								title := strings.TrimSpace()
								link := strings.TrimSpace()

				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: helper.AddFeedbinPage(link),
					Author:      "",
					PubDate:     "",
				}
			}).ToList(&resp)
			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}`, map[string]interface{}{
		"PkgPath": pkgName,
		"Link":    url,
	})
}

func genTest(path, pkgPath string) (string, error) {
	return helper.BuildTemplate(`package {{ .PkgPath }}_test

import (
	"testing"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/{{.Path}}"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_Feed(t *testing.T) {
	as := assert.New(t)

	feed, err := fetch.Fetch({{ .PkgPath }}.New, nil)
	as.Nil(err)

	as.Equal("", feed.Title)
	as.Equal("", feed.Link)
	spew.Dump(feed)
	as.True(len(feed.Items) > 0)
	as.NotEqual(feed.Items[0].Link, feed.Items[1].Link)
	as.NotEmpty(feed.Items[0].Description)
}
`, map[string]interface{}{
		"PkgPath": pkgPath,
		"Path":    path,
	})
}
