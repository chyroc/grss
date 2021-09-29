package zhihu_bookstore_newest

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

func New(map[string]string) (*fetch.Source, error) {
	return &fetch.Source{
		Title: "知乎书店 - 新书抢鲜",
		Link:  "https://www.zhihu.com/pub/features/new",

		Fetch: func() (interface{}, error) {
			resp := new(zhihuBookstoreResp)
			return resp, helper.Req.New(http.MethodGet, "https://api.zhihu.com/books/features/new").Unmarshal(resp)
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj).Transfer(func(obj interface{}) interface{} {
				return obj.(*zhihuBookstoreResp).Data
			}).MapList(func(idx int, obj interface{}) interface{} {
				item := obj.(*zhihuBookstoreRespItem)

				authers, _ := lambda.New(item).Transfer(func(obj interface{}) interface{} {
					return obj.(*zhihuBookstoreRespItem).Authors
				}).MapList(func(idx int, v interface{}) interface{} {
					return v.(*zhihuBookstoreRespItemAuther).Name
				}).ToJoin("、")

				img := regexp.MustCompile(`_.+\.jpg`).ReplaceAllString(item.Cover, ".jpg")

				return &fetch.Item{
					Title: item.Title,
					Link:  item.URL,
					Description: fmt.Sprintf(`<img src="%s"><br>
          <strong>%s</strong><br>
          作者: %s<br><br>
          %s<br><br>
          价格: %d元`, img, item.Title, authers, item.Description, item.Promotion.Price/100),
				}
			}).ToObject(&resp)
			return resp, err
		},
	}, nil
}

type zhihuBookstoreResp struct {
	Data []*zhihuBookstoreRespItem `json:"data"`
}

type zhihuBookstoreRespItemAuther struct {
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Gender    int    `json:"gender"`
	Type      string `json:"type"`
	ID        string `json:"id"`
}

type zhihuBookstoreRespItem struct {
	SkuID       string                          `json:"sku_id"`
	Type        string                          `json:"type"`
	Description string                          `json:"description"`
	Title       string                          `json:"title"`
	URL         string                          `json:"url"`
	Cover       string                          `json:"cover"`
	BookSize    int                             `json:"book_size"`
	ID          int                             `json:"id"`
	IsOwn       bool                            `json:"is_own"`
	Authors     []*zhihuBookstoreRespItemAuther `json:"authors"`
	BookVersion string                          `json:"book_version"`
	Score       float64                         `json:"score"`
	CornerText  string                          `json:"corner_text"`
	Promotion   struct {
		PayType     string `json:"pay_type"`
		IsPromotion bool   `json:"is_promotion"`
		ZhihuBean   int    `json:"zhihu_bean"`
		Price       int    `json:"price"`
		OriginPrice int    `json:"origin_price"`
	} `json:"promotion"`
	BookHash string `json:"book_hash"`
}
