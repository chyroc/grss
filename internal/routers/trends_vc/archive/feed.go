package trends_vc_archive

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://trends.vc/archives/"
	return &fetch.Source{
		Title:       "Trends.vc",
		Description: "Markets and Ideas",
		Link:        link,

		Fetch: func() (interface{}, error) {
			text, err := helper.Req.New(http.MethodGet, link).Text()
			return text, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(string)))
			if err != nil {
				return nil, err
			}

			err = lambda.New(helper.Selection2List(doc.Find("div.entry-content > ul > li"))).MapArrayAsync(func(idx int, obj interface{}) interface{} {
				s := obj.(*goquery.Selection)
				title := strings.TrimSpace(s.Find("a").Text())
				link := strings.TrimSpace(s.Find("a").AttrOr("href", ""))
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: helper.AddFeedbinPage(link),
				}
			}).ToObject(&resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	}, nil
}
