package toutiaoio_index

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://toutiao.io/"
	return &fetch.Source{
		Title:       "开发者头条 - 首页",
		Description: "开发者头条 - 首页",
		Link:        link,

		Fetch: func() (interface{}, error) {
			return nil, nil
		},
		Parse: func(_ interface{}) ([]*fetch.Item, error) {
			all := []*fetch.Item{}
			now := time.Now()
			urls := []string{
				fmt.Sprintf("https://toutiao.io/prev/%s", now.Format("2006-01-02")),
				fmt.Sprintf("https://toutiao.io/prev/%s", now.Add(-time.Hour*24).Format("2006-01-02")),
				fmt.Sprintf("https://toutiao.io/prev/%s", now.Add(-time.Hour*48).Format("2006-01-02")),
			}
			err := lambda.New(urls).MapArrayAsyncWithErr(func(idx int, obj interface{}) (interface{}, error) {
				text, err := helper.Req.New(http.MethodGet, obj.(string)).Text()
				if err != nil {
					return nil, err
				}
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
				if err != nil {
					return nil, err
				}

				sels := []*goquery.Selection{}
				doc.Find(".posts > div").Each(func(i int, sel *goquery.Selection) {
					sels = append(sels, sel)
				})

				items := []*fetch.Item{}
				err = lambda.New(sels).MapArrayAsync(func(idx int, v interface{}) interface{} {
					a := v.(*goquery.Selection).Find(".content > h3 > a")
					title := strings.TrimSpace(a.Text())
					link := strings.TrimSpace(a.AttrOr("href", ""))
					if title == "" || link == "" {
						return nil
					}
					link = "https://toutiao.io" + link
					return &fetch.Item{
						Title:       title,
						Link:        link,
						Description: helper.AddFeedbinPage(link),
					}
				}).FilterList(func(idx int, obj interface{}) bool {
					return obj != nil && obj.(*fetch.Item) != nil
				}).ToObject(&items)
				return items, err
			}).Flatten().ToObject(&all)
			return all, err
		},
	}, nil
}
