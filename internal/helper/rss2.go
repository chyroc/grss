package helper

import (
	"encoding/xml"
	"fmt"

	"github.com/codesoap/rss2"
)

// Create a new rss channel.
func NewChannel(title, link, description string) (*rss2.Channel, error) {
	if len(title) == 0 || len(link) == 0 {
		return nil, fmt.Errorf(`empty string passed to NewChannel()`)
	}
	return &rss2.Channel{
		XMLName:     xml.Name{Local: `channel`},
		Title:       title,
		Link:        link,
		Description: description,
	}, nil
}

// Create a new rss item. Either title or description may be empty.
func NewItem(title, description string) (*rss2.Item, error) {
	if len(title) == 0 {
		return nil, fmt.Errorf(`cannot create item with empty title and description`)
	}
	return &rss2.Item{
		XMLName:     xml.Name{Local: `item`},
		Title:       title,
		Description: description,
	}, nil
}
