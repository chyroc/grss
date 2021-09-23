package grss

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/codesoap/rss2"
)

func DumpFeed(path string, feed *fetch.Feed) error {
	// config
	path = strings.TrimLeft(path, "/")
	date := time.Now().Format("2006-01-02")

	jsonFile := fmt.Sprintf("json/%s/%s.json", date, path)
	xmlFile := fmt.Sprintf("xml/%s.xml", path)

	// load exist data
	oldFeed, err := loadExistFeed(jsonFile)
	if err != nil {
		return err
	}

	// join date
	newFeed, changed := joinFeed(oldFeed, feed)
	if !changed {
		log.Printf("join feed, old=%d, new=%d, no changed, skip", len(oldFeed.Items), len(newFeed.Items))
		return nil
	}

	// save file
	if err = saveJson(jsonFile, newFeed); err != nil {
		return err
	}
	return saveXml(xmlFile, newFeed)
}

// load exist data
func loadExistFeed(jsonFile string) (*fetch.Feed, error) {
	oldFeed := new(fetch.Feed)
	bs, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return oldFeed, nil
	}
	err = json.Unmarshal(bs, oldFeed)
	return oldFeed, err
}

func joinFeed(oldFeed, newFeed *fetch.Feed) (*fetch.Feed, bool) {
	changed := false
	done := map[string]bool{}
	for _, v := range newFeed.Items {
		done[v.Link] = true
	}
	for _, v := range oldFeed.Items {
		if !done[v.Link] {
			done[v.Link] = true
			newFeed.Items = append(newFeed.Items, v)
			changed = true
		}
	}
	return newFeed, changed || len(oldFeed.Items) == 0 && len(newFeed.Items) > 0
}

func saveJson(jsonFile string, feed *fetch.Feed) error {
	if err := os.MkdirAll(filepath.Dir(jsonFile), 0o777); err != nil {
		return err
	}
	bs, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(jsonFile, bs, 0o666)
}

func saveXml(xmlFile string, feed *fetch.Feed) error {
	if err := os.MkdirAll(filepath.Dir(xmlFile), 0o777); err != nil {
		return err
	}
	ch := &rss2.Channel{
		XMLName:     xml.Name{Local: `channel`},
		Title:       feed.Title,
		Link:        feed.Link,
		Description: feed.Description,
		Items:       nil,
	}
	for _, v := range feed.Items {
		item := &rss2.Item{
			XMLName:     xml.Name{Local: `item`},
			Title:       v.Title,
			Link:        v.Link,
			Author:      v.Author,
			Description: v.Description,
		}
		if !v.PubDate.IsZero() {
			item.PubDate = &rss2.RSSTime{Time: v.PubDate}
		}
		ch.Items = append(ch.Items, item)
	}
	rss := rss2.NewRSS(ch)
	bs, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(xmlFile, bs, 0o666)
}
