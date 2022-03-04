package helper

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chyroc/go-feedbin"
)

var Feedbin = feedbin.New(
	feedbin.WithCredential(os.Getenv("FEEDBIN_USERNAME"), os.Getenv("FEEDBIN_PASSWORD")),
	feedbin.WithTimeout(time.Second*10),
)

func AddFeedbinPage(url string) string {
	if url == "" {
		log.Printf("[feedbin] create_page empty url")
		return ""
	}
	if IsInCI {
		text, _ := Req.New(http.MethodGet, url).Text()
		return text
	}

	for i := 0; i < 3; i++ {
		resp, err := Feedbin.CreatePage(context.Background(), &feedbin.CreatePageReq{URL: url})
		if err != nil {
			log.Printf("[feedbin] create_page %q failed: %s", url, err)
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
