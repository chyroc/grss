package xueqiu_snb_article

import (
	"fmt"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	"github.com/chyroc/grss/internal/routers/xueqiu/internal"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://xueqiu.com/?category=snb_article"
	return &fetch.Source{
		Title: "雪球 - 热帖",
		Link:  link,

		Fetch: func() (interface{}, error) {
			uri := fmt.Sprintf("https://xueqiu.com/statuses/hot/listV2.json?since_id=-1&max_id=-1&size=20")
			resp := new(xueqiuSnbArticleResp)
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
			err = lambda.New(obj.(*xueqiuSnbArticleResp).Items).MapList(func(idx int, obj interface{}) interface{} {
				item := obj.(*xueqiuSnbArticleRespItem)
				return &fetch.Item{
					Title:       item.Title(),
					Link:        fmt.Sprintf("https://xueqiu.com%s", item.OriginalStatus.Target),
					Description: item.OriginalStatus.Text,
					Author:      item.OriginalStatus.User.ScreenName,
					PubDate:     time.Unix(item.OriginalStatus.CreatedAt/1000, 0),
				}
			}).ToObject(&resp)

			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}

type xueqiuSnbArticleResp struct {
	Err       string                      `json:"error_description"`
	NextMaxID int                         `json:"next_max_id"`
	Items     []*xueqiuSnbArticleRespItem `json:"items"`
	NextID    int                         `json:"next_id"`
}

type xueqiuSnbArticleRespItem struct {
	ID             int `json:"id"`
	OriginalStatus struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`       // 空
		Description string `json:"description"` // 描述
		UserID      int    `json:"user_id"`
		CreatedAt   int64  `json:"created_at"`
		User        struct {
			ID         int    `json:"id"`
			ScreenName string `json:"screen_name"`
		} `json:"user"`
		Target string `json:"target"`
		Text   string `json:"text"`
	} `json:"original_status"`
}

func (r *xueqiuSnbArticleRespItem) Title() string {
	if r.OriginalStatus.Title != "" {
		return r.OriginalStatus.Title
	}
	return helper.ToTitleText(r.OriginalStatus.Description, 100, " ...")
}
