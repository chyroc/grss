package weibo_user_origin

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

// args: map[string]string{"uid": "1088413295"}
func New(args map[string]string) (*fetch.Source, error) {
	uid := args["uid"]

	userInfo, err := getUserInfo(uid)
	if err != nil {
		return nil, err
	}
	containerID := helper.GetOneMatchString(userInfo.Data.More, regexp.MustCompile(`(\d+)`))

	return &fetch.Source{
		Title:       fmt.Sprintf("微博 - %s 原创微博", userInfo.Data.User.ScreenName),
		Description: userInfo.Data.User.Description,
		Link:        fmt.Sprintf("https://m.weibo.cn/p/%s", uid),

		Fetch: func() (interface{}, error) {
			cards := []*getContainerRespCard{}
			err := lambda.New([]int{1, 2}).MapArrayAsyncWithErr(func(idx int, obj interface{}) (interface{}, error) {
				cards, err := getContainerCard(uid, containerID, obj.(int))
				return cards, err
			}).Flatten().ToList(&cards)
			return cards, err
		},
		Parse: func(obj interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(obj.([]*getContainerRespCard)).FilterArray(func(idx int, obj interface{}) bool {
				return obj.(*getContainerRespCard).Mblog != nil && obj.(*getContainerRespCard).Mblog.RetweetedStatus == nil
			}).MapArray(func(idx int, obj interface{}) interface{} {
				item := obj.(*getContainerRespCard)

				title := strings.TrimSpace(helper.ToTitleText(item.Mblog.Text, 100, " ..."))
				link := fmt.Sprintf("https://m.weibo.cn/detail/%s", item.Mblog.ID)
				desc := item.Mblog.Text
				for _, pic := range item.Mblog.Pics {
					desc += fmt.Sprintf("<div><img src=%q/></div>", pic.Large.URL)
				}
				pubTime, _ := time.Parse("Mon Jan 02 15:04:05 -0700 2006", item.Mblog.CreatedAt)

				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: desc,
					Author:      item.Mblog.User.ScreenName,
					PubDate:     pubTime,
				}
			}).ToList(&resp)
			if err != nil {
				return nil, err
			}
			return resp, err
		},
	}, nil
}

func getUserInfo(uid string) (*userInfo, error) {
	uri := fmt.Sprintf("https://m.weibo.cn/profile/info?uid=%s", uid)
	resp := new(userInfo)
	err := helper.Req.New(http.MethodGet, uri).Unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type userInfo struct {
	Data struct {
		User struct {
			ScreenName      string `json:"screen_name"`
			Description     string `json:"description"`
			ProfileImageURL string `json:"profile_image_url"`
		}
		More string
	}
}

func getContainerCard(uid, containerID string, page int) ([]*getContainerRespCard, error) {
	uri := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?uid=%s&containerid=%s&since_id=0&page=%d", uid, containerID, page)
	resp := new(getContainerResp)
	err := helper.Req.New(http.MethodGet, uri).WithHeaders(map[string]string{
		"Referer":          "https://m.weibo.cn/u/" + uid,
		"MWeibo-Pwa":       "1",
		"X-Requested-With": "XMLHttpRequest",
	}).Unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return resp.Data.Cards, nil
}

type getContainerResp struct {
	Data struct {
		Cards []*getContainerRespCard `json:"cards"`
	} `json:"data"`
}

type getContainerRespCard struct {
	Mblog *getContainerRespCardMblog `json:"mblog"`
}

type getContainerRespCardMblog struct {
	RetweetedStatus interface{} `json:"retweeted_status"`
	// "created_at": "Sun Sep 26 23:57:09 +0800 2021",
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Mid       string `json:"mid"`
	Text      string `json:"text"`
	User      struct {
		ID         int    `json:"id"`
		ScreenName string `json:"screen_name"`
	} `json:"user"`
	IsLongText bool `json:"isLongText"`
	PicNum     int  `json:"pic_num"`
	Pics       []struct {
		Pid  string `json:"pid"`
		URL  string `json:"url"`
		Size string `json:"size"`
		Geo  struct {
			Width  int  `json:"width"`
			Height int  `json:"height"`
			Croped bool `json:"croped"`
		} `json:"geo"`
		Large struct {
			Size string `json:"size"`
			URL  string `json:"url"`
			Geo  struct {
				Width  string `json:"width"`
				Height string `json:"height"`
				Croped bool   `json:"croped"`
			} `json:"geo"`
		} `json:"large"`
	} `json:"pics"`
	Bid string `json:"bid"`
}
