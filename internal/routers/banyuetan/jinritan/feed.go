package banyuetan_jinritan

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "http://www.banyuetan.org/byt/jinritan/index.html"
	return &fetch.Source{
		Title:       "半月谈 - 今日谈",
		Description: "半月谈 - 今日谈",
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
			doc.Find(".byt_tbtj_content").Each(func(i int, sel *goquery.Selection) {
				sels = append(sels, sel)
			})

			items := []*fetch.Item{}
			err = lambda.New(sels).MapArrayAsync(func(idx int, v interface{}) interface{} {
				a := v.(*goquery.Selection).Find("a")
				title := strings.TrimSpace(a.Text())
				link := strings.TrimSpace(a.AttrOr("href", ""))
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: helper.AddFeedbinPage(link),
				}
			}).ToList(&items)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	}, nil
}
