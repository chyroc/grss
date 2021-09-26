package helper

import (
	"context"
	"net/http"
	"os"

	"github.com/chyroc/go-feedbin"
)

var Feedbin = feedbin.New(feedbin.WithCredential(os.Getenv("FEEDBIN_USERNAME"), os.Getenv("FEEDBIN_PASSWORD")))

func AddFeedbinPage(url string) (string, error) {
	for i := 0; i < 3; i++ {
		resp, err := Feedbin.CreatePage(context.Background(), &feedbin.CreatePageReq{URL: url})
		if err != nil {
			if i < 3-1 {
				continue
			}
			return "", err
		}
		return resp.Content, nil
	}
	panic("unreachable")
}

func AddFeedbinPage2(url string) string {
	for i := 0; i < 3; i++ {
		resp, err := Feedbin.CreatePage(context.Background(), &feedbin.CreatePageReq{URL: url})
		if err != nil {
			if i < 3-1 {
				continue
			}
			text, _ := Req.New(http.MethodGet, url).Text()
			return text
		}
		return resp.Content
	}
	panic("unreachable")
}
