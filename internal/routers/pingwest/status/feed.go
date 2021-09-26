package pingwest_status

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title:       "品玩 - 实时要闻",
		Description: "品玩 - 实时要闻",
		Link:        "https://www.pingwest.com/status",

		Fetch: func() (interface{}, error) {
			query := map[string]string{"page": "1"}
			header := map[string]string{"Referer": "https://www.pingwest.com"}
			resp := new(pingwestStateResp)
			return resp, helper.Req.New(http.MethodGet, "https://www.pingwest.com/api/state/list").WithQuerys(query).WithHeaders(header).Unmarshal(resp)
		},
		Parse: func(obj interface{}) (items []*fetch.Item, err error) {
			q, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(*pingwestStateResp).Data.List))
			if err != nil {
				return nil, err
			}

			itemSelections := helper.Selection2List(q.Find(`section.item`))
			err = lambda.New(itemSelections).ArrayAsync(func(idx int, obj interface{}) interface{} {
				selection := obj.(*goquery.Selection)

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
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: helper.AddFeedbinPage(link),
					PubDate:     pubDate,
				}
			}).ToList(&items)
			return items, err
		},
	}, nil
}

type pingwestStateResp struct {
	Data struct {
		List string `json:"list"`
	} `json:"data"`
}
