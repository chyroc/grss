package studygolang_index

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	"github.com/gomarkdown/markdown"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://studygolang.com/"
	return &fetch.Source{
		Title:       "Go语言中文网 - 首页",
		Description: "Go语言中文网 - 首页",
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
			doc.Find(".box_white > div").Each(func(i int, sel *goquery.Selection) {
				sels = append(sels, sel)
			})

			items := []*fetch.Item{}
			err = lambda.New(sels).MapArrayAsync(func(idx int, v interface{}) interface{} {
				a := v.(*goquery.Selection).Find("span.item_title > a")
				title := strings.TrimSpace(a.Text())
				link := strings.TrimSpace(a.AttrOr("href", ""))
				if title == "" || link == "" {
					return nil
				}
				link = "https://studygolang.com" + link
				contentText, _ := helper.Req.New(http.MethodGet, link).Text()
				contentDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(contentText))
				md := contentDoc.Find(".markdown-body").Text()
				output := markdown.ToHTML([]byte(md), nil, nil)
				if len(output) == 0 {
					output = []byte(helper.AddFeedbinPage(link))
				}
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: string(output),
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
