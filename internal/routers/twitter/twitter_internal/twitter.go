package twitter_internal

import (
	"sync"

	"github.com/chyroc/gorequests"
	"github.com/chyroc/grss/internal/helper"
)

type Twitter struct {
	guestToken string
	token      string
	once       sync.Once
}

func New() *Twitter {
	return &Twitter{
		token: "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
	}
}

func (r *Twitter) req(method, url string) *gorequests.Request {
	headers := map[string]string{
		"authorization": "Bearer " + r.token,
		"content-type":  "application/json",
	}
	if r.guestToken != "" {
		headers["x-guest-token"] = r.guestToken
	}
	return helper.Req.New(method, url).WithHeaders(headers)
}

var Ins = New()
