package xueqiu_livenews

import (
	"fmt"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/xueqiu/internal"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://xueqiu.com/?category=livenews"
	return &fetch.Source{
		Title: "雪球 - 7X24快讯",
		Link:  link,

		Fetch: func() (interface{}, error) {
			uri := fmt.Sprintf("https://xueqiu.com/statuses/livenews/list.json?since_id=-1&max_id=-1&count=20")
			resp := new(xueqiuLivenewsResp)
			// err := helper.Req.New(http.MethodGet, uri).Unmarshal(resp)
			err := internal.Request(link, uri, resp)
			if err != nil {
				return nil, err
			} else if resp.Err != "" {
				return nil, fmt.Errorf(resp.Err)
			}
			return resp, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj.(*xueqiuLivenewsResp).Items).MapArray(func(idx int, obj interface{}) interface{} {
				item := obj.(*xueqiuLivenewsRespItem)
				return &fetch.Item{
					Title:       item.Title(),
					Link:        internal.JoinURL(item.Target),
					Description: item.Text,
					PubDate:     time.Unix(item.CreatedAt/1000, 0),
				}
			}).ToList(&resp)

			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}

type xueqiuLivenewsResp struct {
	NextMaxID int                       `json:"next_max_id"`
	Items     []*xueqiuLivenewsRespItem `json:"items"`
	NextID    int                       `json:"next_id"`
	Err       string                    `json:"error_description"`
}

type xueqiuLivenewsRespItem struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	Mark       int    `json:"mark"`
	Target     string `json:"target"`
	CreatedAt  int64  `json:"created_at"`
	ViewCount  int    `json:"view_count"`
	StatusID   int    `json:"status_id"`
	ReplyCount int    `json:"reply_count"`
	ShareCount int    `json:"share_count"`
}

func (r *xueqiuLivenewsRespItem) Title() string {
	if len(r.Text) >= 25 {
		return r.Text[:25] + " ..."
	}
	return r.Text
}
