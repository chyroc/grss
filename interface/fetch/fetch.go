package fetch

import (
	"github.com/chyroc/gorequests"
)

type Source struct {
	Method    string
	URL       string
	Title     string
	Link      string
	Resp      interface{}
	MapReduce func(items interface{}) ([]*Item, error)
}

type Item struct {
	Title       string
	Link        string
	Description string
}

type Feed struct {
	Title string
	Link  string
	Items []*Item
}

func Fetch(source Source) (*Feed, error) {
	req := gorequests.New(source.Method, source.URL)
	if source.Resp != nil {
		err := req.Unmarshal(source.Resp)
		if err != nil {
			return nil, err
		}
	}

	items, err := source.MapReduce(source.Resp)
	if err != nil {
		return nil, err
	}
	return &Feed{
		Title: source.Title,
		Link:  source.Link,
		Items: items,
	}, nil
}
