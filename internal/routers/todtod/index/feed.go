package todtod_index

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://2d2d.io/"
	return &fetch.Source{
		Title:       "TO-D 杂志 - 首页",
		Description: "TO-D 杂志 - 首页",
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
			doc.Find(".items-center > div").Each(func(i int, sel *goquery.Selection) {
				sels = append(sels, sel)
			})

			items := []*fetch.Item{}
			err = lambda.New(sels).MapArrayAsync(func(idx int, v interface{}) interface{} {
				a := v.(*goquery.Selection).Find("a")
				title := strings.TrimSpace(a.Find(".font-bold").Text())
				link := strings.TrimSpace(a.AttrOr("href", ""))
				if title == "" || link == "" {
					return nil
				}
				link = "https://2d2d.io" + link
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
