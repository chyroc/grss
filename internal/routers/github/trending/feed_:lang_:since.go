package github_trending

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
)

// args: map[string]string{"lang": "go", "since": "daily"}
// args: map[string]string{"lang": "python", "since": "daily"}
func New(args map[string]string) (*fetch.Source, error) {
	lang := args["lang"]
	since := args["since"]

	return &fetch.Source{
		Title:       fmt.Sprintf("GitHub - Trending - %s - %s", lang, since),
		Description: fmt.Sprintf("GitHub - Trending - %s - %s", lang, since),
		Link:        "https://github.com",

		Fetch: func() (interface{}, error) {
			query := map[string]string{"page": "1"}
			header := map[string]string{"Referer": "https://github.com"}
			url := fmt.Sprintf("https://github.com/trending/%s?since=%s&spoken_language_code=en", lang, since)
			text, err := helper.Req.New(http.MethodGet, url).WithQuerys(query).WithHeaders(header).Text()
			return text, err
		},
		Parse: func(obj interface{}) (items []*fetch.Item, err error) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(obj.(string)))
			if err != nil {
				return nil, err
			}

			var finalErr error
			doc.Find("article[class=Box-row]").Each(func(i int, s *goquery.Selection) {
				if finalErr != nil {
					return
				}
				itemTitle := ""
				description := ""
				language := ""
				star := 0
				link := ""
				{
					title := strings.TrimSpace(s.Find("svg[class*=octicon-repo]").Parent().AttrOr("href", ""))
					titles := strings.Split(title, "/")
					if len(titles) != 3 {
						finalErr = fmt.Errorf("解析标题失败: %s", s.Text())
						return
					}
					itemTitle = fmt.Sprintf("%s / %s", titles[1], titles[2])
					link = fmt.Sprintf("https://github.com/%s/%s", titles[1], titles[2])
					description = strings.TrimSpace(s.Find("p").Text())
					language = strings.TrimSpace(s.Find("span[itemprop=programmingLanguage]").Text())
					starStr := strings.TrimSpace(s.Find(fmt.Sprintf("svg[aria-label=star]")).Parent().Text())
					starStr = strings.ReplaceAll(starStr, ",", "")
					s, _ := strconv.ParseInt(starStr, 10, 64)
					star = int(s)
				}

				items = append(items, &fetch.Item{
					Title: itemTitle,
					Description: fmt.Sprintf(`
<p>
Desc: %s
<br/>
Lang: %s
<br/>
Star: %d
</p>`, description, language, star),
					Link: link,
				})
			})
			if finalErr != nil {
				return nil, finalErr
			}
			return items, err
		},
	}, nil
}
