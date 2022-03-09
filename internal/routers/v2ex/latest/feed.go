package v2ex_latest

import (
	"net/http"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title: "V2EX - 全部主题",
		Link:  "https://www.v2ex.com/?tab=all",

		Fetch: func() (interface{}, error) {
			resp := []*v2exLatestRespItem{}
			err := helper.Req.New(http.MethodGet, "https://www.v2ex.com/api/topics/latest.json").Unmarshal(&resp)
			return resp, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj).MapArrayAsyncWithErr(func(idx int, obj interface{}) (interface{}, error) {
				item := obj.(*v2exLatestRespItem)
				return &fetch.Item{
					Title:       item.Title + " - " + item.Node.Title,
					Link:        item.URL,
					Description: item.ContentRendered,
					Author:      item.Member.Username,
					PubDate:     time.Unix(item.Created, 0),
				}, nil
			}).ToObject(&resp)
			return resp, err
		},
	}, nil
}

type v2exLatestRespItem struct {
	Node struct {
		AvatarLarge      string `json:"avatar_large"`
		Name             string `json:"name"`
		AvatarNormal     string `json:"avatar_normal"`
		Title            string `json:"title"`
		URL              string `json:"url"`
		Topics           int    `json:"topics"`
		Header           string `json:"header"`
		TitleAlternative string `json:"title_alternative"`
		AvatarMini       string `json:"avatar_mini"`
		Stars            int    `json:"stars"`
		ID               int    `json:"id"`
		ParentNodeName   string `json:"parent_node_name"`
	} `json:"node"`
	Member struct {
		Username     string `json:"username"`
		AvatarNormal string `json:"avatar_normal"`
		Bio          string `json:"bio"`
		URL          string `json:"url"`
		Created      int    `json:"created"`
		AvatarLarge  string `json:"avatar_large"`
		AvatarMini   string `json:"avatar_mini"`
		Location     string `json:"location"`
		ID           int    `json:"id"`
	} `json:"member"`
	LastTouched     int    `json:"last_touched"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	Created         int64  `json:"created"`
	Content         string `json:"content"`
	ContentRendered string `json:"content_rendered"`
	LastModified    int    `json:"last_modified"`
	Replies         int    `json:"replies"`
	ID              int    `json:"id"`
}
