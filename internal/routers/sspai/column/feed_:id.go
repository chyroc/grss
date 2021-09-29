package sspai_column

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

// args: map[string]string{"id": "264"}
// args: map[string]string{"id": "266"}
func New(args map[string]string) (*fetch.Source, error) {
	id := args["id"]

	link := fmt.Sprintf("https://sspai.com/column/" + id)
	columnInfoAPI := fmt.Sprintf("https://sspai.com/api/v1/special_columns/" + id)
	columnInfo := new(sspaiColumnResp)
	if err := helper.Req.New(http.MethodGet, columnInfoAPI).WithHeader("Referer", link).Unmarshal(columnInfo); err != nil {
		return nil, err
	}

	return &fetch.Source{
		Title:       "少数派专栏 - " + columnInfo.Title,
		Description: columnInfo.Intro,
		Link:        link,

		Fetch: func() (interface{}, error) {
			url := fmt.Sprintf("https://sspai.com/api/v1/articles?offset=0&limit=10&special_column_ids=%s&include_total=false", id)
			header := map[string]string{"Referer": link}
			resp := new(sspaiColumnItemResp)
			return resp, helper.Req.New(http.MethodGet, url).WithHeaders(header).Unmarshal(resp)
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj).Transfer(func(obj interface{}) interface{} {
				return obj.(*sspaiColumnItemResp).List
			}).MapArrayAsync(func(idx int, v interface{}) interface{} {
				item := v.(*sspaiColumnItemRespItem)
				title := strings.TrimSpace(item.Title)
				pubTime := time.Unix(item.CreatedAt, 0)
				itemURL := fmt.Sprintf("https://sspai.com/post/%d", item.ID)
				author := item.Author.Nickname

				return &fetch.Item{
					Title:       title,
					Link:        itemURL,
					Description: helper.AddFeedbinPage(itemURL),
					Author:      author,
					PubDate:     pubTime,
				}
			}).ToObject(&resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	}, nil
}

type sspaiColumnResp struct {
	ID                  int    `json:"id"`
	CreatedAt           int    `json:"created_at"`
	ReleasedAt          int    `json:"released_at"`
	Title               string `json:"title"`
	Intro               string `json:"intro"`
	Banner              string `json:"banner"`
	BannerID            int    `json:"banner_id"`
	RequestReason       string `json:"request_reason"`
	RejectReason        string `json:"reject_reason"`
	AuthorID            int    `json:"author_id"`
	PendingAuthorID     int    `json:"pending_author_id"`
	Status              string `json:"status"`
	IsMatrix            bool   `json:"is_matrix"`
	AdvertisementID     int    `json:"advertisement_id"`
	ViewsCount          int    `json:"views_count"`
	RecommendedAt       int    `json:"recommended_at"`
	ArticlesCount       int    `json:"articles_count"`
	FollowersCount      int    `json:"followers_count"`
	Followed            bool   `json:"followed"`
	ReleasedOrRetiredAt int    `json:"released_or_retired_at"`
	AcceptArticle       bool   `json:"accept_article"`
	EditorsCount        int    `json:"editors_count"`
	AuthorsCount        int    `json:"authors_count"`
}

type sspaiColumnItemResp struct {
	List []*sspaiColumnItemRespItem `json:"list"`
}

type sspaiColumnItemRespItem struct {
	ID           int    `json:"id"`
	CreatedAt    int64  `json:"created_at"`
	Banner       string `json:"banner"`
	BannerID     int    `json:"banner_id"`
	Title        string `json:"title"`
	ReleasedAt   int    `json:"released_at"`
	ModifyAt     int    `json:"modify_at"`
	Summary      string `json:"summary"`
	WordsCount   int    `json:"words_count"`
	AllowComment bool   `json:"allow_comment"`
	Type         string `json:"type"`
	PostType     int    `json:"post_type"`
	Important    int    `json:"important"`
	Free         int    `json:"free"`
	Author       struct {
		ID           int         `json:"id"`
		Nickname     string      `json:"nickname"`
		Avatar       string      `json:"avatar"`
		AvatarID     int         `json:"avatar_id"`
		Bio          string      `json:"bio"`
		Role         string      `json:"role"`
		EmailMatrixs bool        `json:"email_matrixs"`
		LikedCount   int         `json:"liked_count"`
		Slug         string      `json:"slug"`
		Member       interface{} `json:"member"`
	} `json:"author"`
	Tags []struct {
		ID         int         `json:"id"`
		CreatedAt  int         `json:"created_at"`
		ReleasedAt int         `json:"released_at"`
		ModifyAt   int         `json:"modify_at"`
		Title      string      `json:"title"`
		Intro      string      `json:"intro"`
		ViewsCount int         `json:"views_count"`
		UsableUser bool        `json:"usable_user"`
		Tags       interface{} `json:"tags"`
	} `json:"tags"`
	LikesCount     int `json:"likes_count"`
	SpecialColumns []struct {
		ID            int    `json:"id"`
		CreatedAt     int    `json:"created_at"`
		ReleasedAt    int    `json:"released_at"`
		Title         string `json:"title"`
		Intro         string `json:"intro"`
		Banner        string `json:"banner"`
		BannerID      int    `json:"banner_id"`
		RequestReason string `json:"request_reason"`
		AuthorID      int    `json:"author_id"`
		Status        string `json:"status"`
		ViewsCount    int    `json:"views_count"`
		AuthorsCount  int    `json:"authors_count"`
	} `json:"special_columns"`
	PodcastConfigs interface{} `json:"podcast_configs"`
}
