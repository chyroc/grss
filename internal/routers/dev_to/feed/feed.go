package dev_to_feed

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	link := "https://dev.to"
	return &fetch.Source{
		Title: "DEV Community - Feed",
		Link:  link,

		Fetch: func() (interface{}, error) {
			link := "https://dev.to/search/feed_content?per_page=20&page=0&sort_by=hotness_score&sort_direction=desc&approved=&class_name=Article"
			resp := new(sdevToFeedResp)

			err := helper.Req.New(http.MethodGet, link).Unmarshal(resp)
			if err != nil {
				return nil, err
			}
			return resp.Result, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj.([]*sdevToFeedRespItem)).MapArrayAsync(func(idx int, obj interface{}) interface{} {
				item := obj.(*sdevToFeedRespItem)
				link := fmt.Sprintf("https://dev.to" + item.Path)

				return &fetch.Item{
					Title:       helper.FanyiAndAppend(item.Title, " | "),
					Link:        link,
					Description: helper.FetchFeedBinAndFanyiAndAppend(link),
					Author:      item.User.Username,
					PubDate:     time.Unix(item.PublishedAtInt, 10),
				}
			}).ToObject(&resp)
			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}

type sdevToFeedResp struct {
	Result []*sdevToFeedRespItem `json:"result"`
}

type sdevToFeedRespItem struct {
	ClassName            string      `json:"class_name"`
	CloudinaryVideoURL   interface{} `json:"cloudinary_video_url"`
	CommentsCount        int         `json:"comments_count"`
	ID                   int         `json:"id"`
	Path                 string      `json:"path"`
	PublicReactionsCount int         `json:"public_reactions_count"`
	ReadablePublishDate  string      `json:"readable_publish_date"`
	ReadingTime          int         `json:"reading_time"`
	Title                string      `json:"title"`
	UserID               int         `json:"user_id"`
	VideoDurationString  string      `json:"video_duration_string"`
	PublishedAtInt       int64       `json:"published_at_int"`
	TagList              []string    `json:"tag_list"`
	FlareTag             struct {
		Name         string `json:"name"`
		BgColorHex   string `json:"bg_color_hex"`
		TextColorHex string `json:"text_color_hex"`
	} `json:"flare_tag"`
	User struct {
		Name           string `json:"name"`
		ProfileImage90 string `json:"profile_image_90"`
		Username       string `json:"username"`
	} `json:"user"`
}
