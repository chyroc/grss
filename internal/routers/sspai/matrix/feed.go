package sspai_matrix

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
		Title:       "少数派 - Matrix",
		Description: "少数派 - Matrix",
		Link:        "https://sspai.com/matrix",

		Fetch: func() (interface{}, error) {
			url := "https://sspai.com/api/v1/articles?offset=0&limit=20&is_matrix=1&sort=matrix_at&include_total=false"
			resp := new(sspaiMatrixResp)
			return resp, helper.Req.New(http.MethodGet, url).Unmarshal(resp)
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj).Transfer(func(obj interface{}) interface{} {
				return obj.(*sspaiMatrixResp).List
			}).ArrayAsync(func(idx int, obj interface{}) interface{} {
				item := obj.(*sspaiMatrixRespItem)
				link := fmt.Sprintf("https://sspai.com/post/%d", item.ID)
				// desc := ""
				text, _ := helper.Req.New(http.MethodGet, link).Text()
				// 		// description = $('#app > div.postPage.article-wrapper > div.article-detail > article > div.article-body').html();

				return &fetch.Item{
					Title:       strings.TrimSpace(item.Title),
					Link:        link,
					Description: text,
					Author:      strings.TrimSpace(item.Author.Nickname),
					PubDate:     time.Unix(item.ReleasedAt, 0),
				}
			}).ToList(&resp)
			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}

type sspaiMatrixRespItem struct {
	ID                int    `json:"id"`
	CreatedAt         int    `json:"created_at"`
	Banner            string `json:"banner"`
	BannerID          int    `json:"banner_id"`
	Title             string `json:"title"`
	ReleasedAt        int64  `json:"released_at"`
	ModifyAt          int    `json:"modify_at"`
	Summary           string `json:"summary"`
	WordsCount        int    `json:"words_count"`
	AllowComment      bool   `json:"allow_comment"`
	PromoteTitle      string `json:"promote_title"`
	IsMatrix          bool   `json:"is_matrix"`
	Type              string `json:"type"`
	FollowUpAdminID   int    `json:"follow_up_admin_id"`
	RecommendToHomeAt int    `json:"recommend_to_home_at"`
	MatrixAt          int    `json:"matrix_at"`
	ShowContentTable  bool   `json:"show_content_table"`
	Author            struct {
		ID            int    `json:"id"`
		Nickname      string `json:"nickname"`
		Avatar        string `json:"avatar"`
		AvatarID      int    `json:"avatar_id"`
		Bio           string `json:"bio"`
		Role          string `json:"role"`
		SignedWriter  bool   `json:"signed_writer"`
		EmailMessages bool   `json:"email_messages"`
		LikedCount    int    `json:"liked_count"`
		Slug          string `json:"slug"`
	} `json:"author"`
	Tags []struct {
		ID         int    `json:"id"`
		ReleasedAt int    `json:"released_at"`
		ModifyAt   int    `json:"modify_at"`
		Title      string `json:"title"`
		Intro      string `json:"intro,omitempty"`
		ViewsCount int    `json:"views_count"`
		UsableUser bool   `json:"usable_user"`
	} `json:"tags"`
	LikesCount         int `json:"likes_count"`
	FavoritesCount     int `json:"favorites_count"`
	AllCommentsCount   int `json:"all_comments_count"`
	CommentsCount      int `json:"comments_count"`
	CommentReplysCount int `json:"comment_replys_count"`
	Corners            []struct {
		ID        int    `json:"id"`
		CreatedAt int    `json:"created_at"`
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		Memo      string `json:"memo"`
		Status    int    `json:"status"`
		Color     string `json:"color"`
		Weight    int    `json:"weight"`
		AllowSet  bool   `json:"allow_set"`
	} `json:"corners"`
}

type sspaiMatrixResp struct {
	List []*sspaiMatrixRespItem `json:"list"`
}
