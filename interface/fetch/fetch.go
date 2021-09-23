package fetch

import (
	"fmt"
	"net/url"
	"time"

	"github.com/chyroc/gorequests"
)

type Source struct {
	Title       string
	Description string
	Link        string

	Method    string
	URL       string
	Query     url.Values
	Header    url.Values
	Resp      interface{}
	MapReduce func(obj interface{}) ([]*Item, error)
}

type Item struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
}

type Feed struct {
	Title       string
	Link        string
	Description string
	Items       []*Item
}

func Fetch(source Source) (*Feed, error) {
	req := gorequests.New(source.Method, source.URL)
	if len(source.Query) > 0 {
		for k, v := range source.Query {
			for _, vv := range v {
				req.WithQuery(k, vv)
			}
		}
	}
	if len(source.Header) > 0 {
		for k, v := range source.Header {
			for _, vv := range v {
				req.WithHeader(k, vv)
			}
		}
	}

	if source.Resp != nil {
		err := req.Unmarshal(source.Resp)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println(req.Text())
		panic("")
	}

	items, err := source.MapReduce(source.Resp)
	if err != nil {
		return nil, err
	}
	return &Feed{
		Title:       source.Title,
		Link:        source.Link,
		Description: source.Description,
		Items:       items,
	}, nil
}
