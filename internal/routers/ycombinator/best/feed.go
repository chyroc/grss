package ycombinator_best

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/grss/internal/fetch"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title:  "Hacker News - Top Links",
		Link:   "https://news.ycombinator.com/best",
		Method: http.MethodGet,
		URL:    "https://news.ycombinator.com/best",
		MapReduce: func(obj interface{}) ([]*fetch.Item, error) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(string)))
			if err != nil {
				return nil, err
			}

			items := []*fetch.Item{}
			doc.Find(".athing").Each(func(i int, selection *goquery.Selection) {
				id := selection.AttrOr("id", "")
				link := selection.Find("a.storylink").AttrOr("href", "")
				title := strings.TrimSpace(selection.Find("a.storylink").Text())
				next := selection.Next()
				age := strings.TrimSpace(next.Find(".age").AttrOr("title", ""))
				point := strings.TrimSpace(next.Find(".score").Text())
				pubTime, _ := time.Parse("2006-01-02T15:04:05", age)
				desc := fmt.Sprintf(`Article URL: %s

Comments URL: https://news.ycombinator.com/item?id=%s

Points: %s
`, link, id, point)
				items = append(items, &fetch.Item{
					Title:       title,
					Link:        link,
					Description: desc,
					PubDate:     pubTime,
				})
			})
			return items, nil
		},
	}, nil
}
