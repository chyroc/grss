package post

import (
	"fmt"
	"net/http"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

// args: map[string]string{"r": "chasays"}
// args: map[string]string{"r": "pythonhunter"}
func New(args map[string]string) (*fetch.Source, error) {
	r := args["r"]
	link := fmt.Sprintf("https://%s.zhubai.love/", r)

	postResp := new(postItemResp)
	err := helper.Req.New(http.MethodGet, fmt.Sprintf("https://%s.zhubai.love/api/publications/%s/posts?publication_id_type=token", r, r)).Unmarshal(resp)
	if err != nil {
		return nil, err
	}

	title := "竹白"
	desc := ""

	if len(postResp.Data) > 0 {
		title = postResp.Data[0].Publication.Name
		desc = postResp.Data[0].Publication.Description
	}

	return &fetch.Source{
		Title:       "竹白 - " + title,
		Description: desc,
		Link:        link,
		Fetch: func() (interface{}, error) {
			return nil, nil
		},
		Parse: func(_ interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(postResp.Data).MapArrayAsyncWithErr(func(idx int, obj interface{}) (interface{}, error) {
				item := obj.(*postItem)
				title := item.Title
				link := fmt.Sprintf("https://%s.zhubai.love/posts/%s", r, item.ID)
				text := helper.AddFeedbinPage(link)
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: text,
				}, nil
			}).ToObject(&resp)
			return resp, err
		},
	}, nil
}

type postItem struct {
	Author struct {
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Name        string `json:"name"`
	} `json:"author"`
	Content       string      `json:"content"`
	CreatedAt     int64       `json:"created_at"`
	ID            string      `json:"id"`
	IsPaidContent bool        `json:"is_paid_content"`
	Paywall       interface{} `json:"paywall"`
	Publication   struct {
		CreatedAt   int64       `json:"created_at"`
		Description string      `json:"description"`
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Theme       interface{} `json:"theme"`
		Token       string      `json:"token"`
		UpdatedAt   int64       `json:"updated_at"`
	} `json:"publication"`
	Title     string `json:"title"`
	UpdatedAt int64  `json:"updated_at"`
}

type postItemResp struct {
	Data       []*postItem `json:"data"`
	Pagination struct {
		HasNext bool   `json:"has_next"`
		HasPrev bool   `json:"has_prev"`
		Next    string `json:"next"`
		Prev    string `json:"prev"`
	} `json:"pagination"`
}
