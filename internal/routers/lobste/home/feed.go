package lobste_home

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
	link := "https://lobste.rs/"
	return &fetch.Source{
		Title: "Lobsters - Home",
		Link:  link,

		Fetch: func() (interface{}, error) {
			text, err := helper.Req.New(http.MethodGet, link).Text()
			return text, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(string)))
			if err != nil {
				return nil, err
			}

			err = lambda.New(helper.Selection2List(doc.Find(".story"))).MapArray(func(idx int, obj interface{}) interface{} {
				s := obj.(*goquery.Selection)

				a := s.Find("a.u-url")
				title := strings.TrimSpace(a.Text())
				link := strings.TrimSpace(a.AttrOr("href", ""))
				author := strings.TrimSpace(s.Find("a.u-author").Text())
				pubTime, _ := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(s.Find("div.byline > span").AttrOr("title", "")))
				lrsLink := strings.TrimSpace(s.Find("span.comments_label > a").AttrOr("href", ""))
				if lrsLink != "" {
					lrsLink = "https://lobste.rs" + lrsLink
				}
				linkText := helper.AddFeedbinPage(link)
				lrsLinkText := helper.AddFeedbinPage(lrsLink)
				text := fmt.Sprintf(`<div>Article:<div><br/><br/><div>%s</div><br/><br/><div>Lobsters:<div><br/><br/><div>%s</div>`, linkText, lrsLinkText)
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: text,
					Author:      author,
					PubDate:     pubTime,
				}
			}).ToList(&resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	}, nil
}
