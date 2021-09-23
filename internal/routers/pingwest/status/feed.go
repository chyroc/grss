package pingwest_status

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/grss/internal/fetch"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title:       "品玩 - 实时要闻",
		Description: "品玩 - 实时要闻",
		Link:        "https://www.pingwest.com/status",

		Method: http.MethodGet,
		URL:    "https://www.pingwest.com/api/state/list",
		Query:  map[string][]string{"page": {"1"}},
		Header: map[string][]string{"Referer": {"https://www.pingwest.com"}},
		Resp:   new(pingwestStateResp),
		MapReduce: func(obj interface{}) (items []*fetch.Item, err error) {
			q, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(*pingwestStateResp).Data.List))
			if err != nil {
				return nil, err
			}

			q.Find(`section.item`).Each(func(i int, selection *goquery.Selection) {
				timestamp := strings.TrimSpace(selection.AttrOr("data-t", ""))
				ts, _ := strconv.ParseInt(timestamp, 10, 64)
				pubDate := time.Now()
				if ts > 0 {
					pubDate = time.Unix(ts, 0)
				}
				rightNode := selection.Find(`.news-info`)
				tag := strings.TrimSpace(rightNode.Find(".item-tag-list").Text())
				title := strings.TrimSpace(rightNode.Find(".title").Text())
				if title == "" {
					title = tag
				}
				link := strings.TrimSpace(rightNode.Find("a").Last().AttrOr("href", ""))
				if !strings.HasPrefix(link, "http") {
					link = "https:" + link
				}
				description := strings.TrimSpace(rightNode.Text())
				items = append(items, &fetch.Item{
					Title:       title,
					Link:        link,
					Description: description,
					PubDate:     pubDate,
				})
			})

			return items, err
		},
	}, nil
}

type pingwestStateResp struct {
	Data struct {
		List string `json:"list"`
	} `json:"data"`
}
