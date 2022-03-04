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
	"github.com/chyroc/grss/internal/helper"
	"github.com/codesoap/rss2"
)

func DumpFeed(path string, feed *fetch.Feed) error {
	// config
	path = strings.TrimLeft(path, "/")

	// refresh data && save history data
	{
		date := time.Now().Format("2006-01-02")
		jsonDateFile := fmt.Sprintf("json/%s/%s.json", date, path)
		oldFeed, err := loadExistFeed(jsonDateFile)
		if err != nil {
			return err
		}

		// ToJoin date
		newFeed, changed := ToJoinFeed(oldFeed, feed)
		if !changed {
			log.Printf("ToJoin feed, old=%d, new=%d, no changed, skip", len(oldFeed.Items), len(newFeed.Items))
			return nil
		}
		if err = saveJson(jsonDateFile, newFeed); err != nil {
			return err
		}

		feed = newFeed
	}

	if len(feed.Items) == 0 {
		return nil
	}

	// save the latest file
	{
		jsonFile := fmt.Sprintf("json/%s/%s.json", "latest", path)
		xmlFile := fmt.Sprintf("xml/%s.xml", path)

		if err := saveJson(jsonFile, feed); err != nil {
			return err
		}
		if err := saveXml(xmlFile, feed); err != nil {
			return err
		}
	}

	return nil
}

func removeOldData() {
	date := time.Now()
	for i := 3; i < 10; i++ {
		oldDate := date.AddDate(0, 0, -i).Format("2006-01-02")
		jsonDir := fmt.Sprintf("json/%s", oldDate)
		_ = os.Remove(jsonDir)
	}
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

func ToJoinFeed(oldFeed, newFeed *fetch.Feed) (*fetch.Feed, bool) {
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
	ch, err := helper.NewChannel(feed.Title, feed.Link, feed.Description)
	if err != nil {
		return err
	}
	for _, v := range feed.Items {
		item, err := helper.NewItem(v.Title, v.Description)
		if err != nil {
			return err
		}
		item.Link = v.Link
		item.Author = v.Author
		if !(v.PubDate.IsZero() || v.PubDate.Year() == 0 || v.PubDate.Year() == 1) {
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
