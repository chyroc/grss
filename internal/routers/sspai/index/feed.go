package sspai_index

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title:       "少数派 - 推荐",
		Description: "少数派 - 推荐",
		Link:        "https://sspai.com",

		Fetch: func() (interface{}, error) {
			url := "https://sspai.com/api/v1/article/index/page/get?limit=30&offset=0"
			resp := new(sspaiResp)
			return resp, helper.Req.New(http.MethodGet, url).Unmarshal(resp)
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj).Transfer(func(obj interface{}) interface{} {
				return obj.(*sspaiResp).Data
			}).MapArrayAsync(func(idx int, obj interface{}) interface{} {
				item := obj.(*sspaiItem)
				link := fmt.Sprintf("https://sspai.com/post/%d", item.ID)

				return &fetch.Item{
					Title:       strings.TrimSpace(item.Title),
					Link:        link,
					Description: helper.AddFeedbinPage(link),
					Author:      strings.TrimSpace(item.Author.Nickname),
					PubDate:     time.Unix(item.ReleasedTime, 0),
				}
			}).ToObject(&resp)
			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}

type sspaiItem struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	ReleasedTime int64  `json:"released_time"`
	Author       struct {
		ID       int    `json:"id"`
		Slug     string `json:"slug"`
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	} `json:"author"`
	CreatedTime int `json:"created_time"`
	ModifyTime  int `json:"modify_time"`
}

type sspaiResp struct {
	Error int          `json:"error"`
	Msg   string       `json:"msg"`
	Data  []*sspaiItem `json:"data"`
	Total int          `json:"total"`
}
