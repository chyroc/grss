package meituan_tech_article

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://tech.meituan.com/"
	return &fetch.Source{
		Title:       "美团技术团队 - 文章",
		Description: "美团技术团队 - 文章",
		Link:        link,

		Fetch: func() (interface{}, error) {
			text, err := helper.Req.New(http.MethodGet, link).Text()
			return text, err
		},
		Parse: func(obj interface{}) ([]*fetch.Item, error) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(string)))
			if err != nil {
				return nil, err
			}

			sels := []*goquery.Selection{}
			doc.Find(".post-container").Each(func(i int, sel *goquery.Selection) {
				sels = append(sels, sel)
			})

			items := []*fetch.Item{}
			err = lambda.New(sels).MapArrayAsync(func(idx int, v interface{}) interface{} {
				a := v.(*goquery.Selection).Find(".post-title > a")
				title := strings.TrimSpace(a.Text())
				link := strings.TrimSpace(a.AttrOr("href", ""))
				if title == "" || link == "" {
					return nil
				}
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: helper.AddFeedbinPage(link),
				}
			}).FilterList(func(idx int, obj interface{}) bool {
				return obj != nil && obj.(*fetch.Item) != nil
			}).ToObject(&items)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	}, nil
}
