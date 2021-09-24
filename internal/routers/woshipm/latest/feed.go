package woshipm_latest

import (
	"net/http"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title:       "人人都是产品经理 - 最新文章",
		Description: "",
		Link:        "http://www.woshipm.com",

		Fetch: func() (interface{}, error) {
			items := []*woshipmLatestRespItem{}
			f := lambda.New([]string{
				"http://www.woshipm.com/__api/v1/stream-list",
				"http://www.woshipm.com/__api/v1/stream-list/page/2",
				"http://www.woshipm.com/__api/v1/stream-list/page/3",
			})
			err := f.ArrayAsyncWithErr(func(idx int, v interface{}) (interface{}, error) {
				resp := new(woshipmLatestResp)
				err := helper.Req.New(http.MethodGet, v.(string)).Unmarshal(resp)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}).Flatten().ToList(&items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
		Parse: func(obj interface{}) ([]*fetch.Item, error) {
			items:=[]*fetch.Item{}
			err:=lambda.New(obj).ArrayAsync(func(idx int, obj interface{}) interface{} {
				item := obj.(*woshipmLatestRespItem)
				title := item.Title
				link := item.Permalink
				pubDate, _ := time.Parse("2006/01/02", item.Date)
				text, _ := helper.Req.New(http.MethodGet, link).Text()

				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: text,
					PubDate:     pubDate,
				}
			}).ToList(&items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	}, nil
}

type woshipmLatestResp struct {
	Success string                   `json:"success"`
	Payload []*woshipmLatestRespItem `json:"payload"`
}

type woshipmLatestRespItem struct {
	ID        int    `json:"id"`
	IsEvent   bool   `json:"is_event"`
	PostType  string `json:"post_type"`
	Title     string `json:"title"`
	Permalink string `json:"permalink"`
	Date      string `json:"date"`
	Author    struct {
		Name   string `json:"name"`
		ID     int    `json:"id"`
		Avatar string `json:"avatar"`
		Link   string `json:"link"`
		Role   string `json:"role"`
	} `json:"author"`
	Image    string `json:"image"`
	Category string `json:"category"`
	Snipper  string `json:"snipper"`
	// Like      interface{}   `json:"like"`
	// View      interface{}   `json:"view"`
	// Bookmark  interface{}   `json:"bookmark"`
	// EventInfo []interface{} `json:"event_info"`
	Comment string `json:"comment"`
	Catslug string `json:"catslug"`
}
