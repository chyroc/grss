package reddit_community_hot

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/baidufanyi"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

// {"r": "golang"}
func New(args map[string]string) (*fetch.Source, error) {
	r := args["r"]
	link := fmt.Sprintf("https://www.reddit.com/r/%s/hot/", r)

	text, err := helper.Req.New(http.MethodGet, link).Text()
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	desc := strings.TrimSpace(doc.Find("meta[name=description]").AttrOr("content", ""))

	return &fetch.Source{
		Title:       "Reddit - " + title,
		Description: desc,
		Link:        link,

		Fetch: func() (interface{}, error) {
			return doc, nil
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			doc := obj.(*goquery.Document)
			containers := []*goquery.Selection{}
			doc.Find("div[data-testid=post-container]").Each(func(i int, selection *goquery.Selection) {
				containers = append(containers, selection)
			})
			err = lambda.New(containers).ArrayAsyncWithErr(func(idx int, obj interface{}) (interface{}, error) {
				selection := obj.(*goquery.Selection)
				a := selection.Find("a[data-click-id=body]")
				link := strings.TrimSpace(a.AttrOr("href", ""))
				if link != "" {
					link = "https://www.reddit.com" + link
				}
				text, _ := helper.Req.New(http.MethodGet, link).Text()
				title := strings.TrimSpace(a.Find("h3").Text())
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: text,
				}, nil
			}).ToList(&resp)
			return resp, err
		},
	}, nil
}

func init() {
	_ = baidufanyi.New(baidufanyi.WithCredential("", ""))
}
