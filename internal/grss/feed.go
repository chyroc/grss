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
		htmlDir := fmt.Sprintf("html/%s", path)

		if err := saveJson(jsonFile, feed); err != nil {
			return err
		}
		if err := saveXml(xmlFile, feed); err != nil {
			return err
		}
		if err := saveHtml(htmlDir, feed); err != nil {
			return err
		}
	}

	return nil
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

func saveHtml(htmlDir string, feed *fetch.Feed) error {
	// 删除之前运行产生的 html
	if len(feed.Items) > 0 {
		fs, err := ioutil.ReadDir(htmlDir)
		if err != nil {
			return err
		}
		for _, v := range fs {
			if strings.HasSuffix(v.Name(), ".html") {
				_ = os.Remove(htmlDir + "/" + v.Name())
			}
		}
	}

	// 生成新的 html
	if err := os.MkdirAll(htmlDir, 0o777); err != nil {
		return err
	}
	{
		indexMdContent, err := generateFeedHtml(feed)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(fmt.Sprintf("%s/index.html", htmlDir), []byte(indexMdContent), 0o666); err != nil {
			return err
		}
		for _, v := range feed.Items {
			if err = ioutil.WriteFile(fmt.Sprintf("%s/%s.html", htmlDir, helper.Md5(v.Link)), []byte(v.Description), 0o666); err != nil {
				return err
			}
		}
	}

	return nil
}

func generateFeedHtml(feed *fetch.Feed) (string, error) {
	type Item struct {
		Title   string
		Link    string
		LinkMd5 string
	}

	type Feed struct {
		Title string
		Link  string
		Items []*Item
	}
	data := &Feed{Title: feed.Title, Link: feed.Link}
	for _, v := range feed.Items {
		data.Items = append(data.Items, &Item{Title: v.Title, Link: v.Link, LinkMd5: helper.Md5(v.Link)})
	}

	return helper.BuildTemplate(`<!DOCTYPE html><html>

<head>
	<title>{{ .Title }}</title>
</head>

<body>
	<ul>
{{ range .Items }}
		<li><a href="./{{ .LinkMd5 }}.html">{{ .Title }}</a></li>
{{ end }}
	</ul>
</body>

</html>
`, data)
}
