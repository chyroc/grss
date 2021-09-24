package fetch

import (
	"time"
)

type Source struct {
	Title       string
	Description string
	Link        string

	Fetch func() (interface{}, error)
	Parse func(obj interface{}) ([]*Item, error)
}

type Item struct {
	Title       string
	Link        string
	Description string
	Author      string
	PubDate     time.Time
}

type Feed struct {
	Title       string
	Link        string
	Description string
	Items       []*Item
}

func Fetch(sourceGetter func(map[string]string) (*Source, error), args map[string]string) (*Feed, error) {
	source, err := sourceGetter(args)
	if err != nil {
		return nil, err
	}

	resp, err := source.Fetch()
	if err != nil {
		return nil, err
	}

	items, err := source.Parse(resp)
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
