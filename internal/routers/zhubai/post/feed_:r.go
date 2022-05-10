package zhubai_post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chyroc/go-lambda"
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	"github.com/gomarkdown/markdown"
)

// args: map[string]string{"r": "chasays"}
// args: map[string]string{"r": "pythonhunter"}
// args: map[string]string{"r": "produck"}
// args: map[string]string{"r": "oistore"}
// args: map[string]string{"r": "decohack"}
func New(args map[string]string) (*fetch.Source, error) {
	r := args["r"]
	link := fmt.Sprintf("https://%s.zhubai.love/", r)

	postResp := new(postItemResp)
	err := helper.Req.New(http.MethodGet, fmt.Sprintf("https://%s.zhubai.love/api/publications/%s/posts?publication_id_type=token", r, r)).Unmarshal(postResp)
	if err != nil {
		return nil, err
	}

	title := "竹白"
	desc := ""

	if len(postResp.Data) > 0 {
		title = postResp.Data[0].Publication.Name
		desc = postResp.Data[0].Publication.Description
	}

	return &fetch.Source{
		Title:       "竹白 - " + title,
		Description: desc,
		Link:        link,
		Fetch: func() (interface{}, error) {
			return nil, nil
		},
		Parse: func(_ interface{}) (resp []*fetch.Item, err error) {
			err = lambda.New(postResp.Data).MapArrayAsyncWithErr(func(idx int, obj interface{}) (interface{}, error) {
				item := obj.(*postItem)
				title := item.Title
				link := fmt.Sprintf("https://%s.zhubai.love/posts/%s", r, item.ID)
				blocks := []*ZhubaiBlock{}
				_ = json.Unmarshal([]byte(item.Content), &blocks)
				output := markdown.ToHTML([]byte(zhubaiBlockListMarkdown(blocks)), nil, nil)
				return &fetch.Item{
					Title:       title,
					Link:        link,
					Description: string(output),
				}, nil
			}).ToObject(&resp)
			return resp, err
		},
	}, nil
}

type postItem struct {
	Author struct {
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Name        string `json:"name"`
	} `json:"author"`
	Content       string      `json:"content_json"`
	CreatedAt     int64       `json:"created_at"`
	ID            string      `json:"id"`
	IsPaidContent bool        `json:"is_paid_content"`
	Paywall       interface{} `json:"paywall"`
	Publication   struct {
		CreatedAt   int64       `json:"created_at"`
		Description string      `json:"description"`
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Theme       interface{} `json:"theme"`
		Token       string      `json:"token"`
		UpdatedAt   int64       `json:"updated_at"`
	} `json:"publication"`
	Title     string `json:"title"`
	UpdatedAt int64  `json:"updated_at"`
}

type postItemResp struct {
	Data       []*postItem `json:"data"`
	Pagination struct {
		HasNext bool   `json:"has_next"`
		HasPrev bool   `json:"has_prev"`
		Next    string `json:"next"`
		Prev    string `json:"prev"`
	} `json:"pagination"`
}

type ZhubaiBlock struct {
	Type     string         `json:"type"` // paragraph,bulleted-list
	Children []*ZhubaiBlock `json:"children"`

	URL    string `json:"url,omitempty"`
	Bold   bool   `json:"bold"`
	Italic bool   `json:"italic"`
	Text   string `json:"text"`
}

func zhubaiBlockListMarkdown(blocks []*ZhubaiBlock) string {
	r := new(strings.Builder)
	for _, block := range blocks {
		r.WriteString(block.zhubaiBlockMarkdown())
	}
	return r.String()
}

func (block *ZhubaiBlock) zhubaiBlockMarkdown() string {
	r := new(strings.Builder)
	switch block.Type {
	case "paragraph":
		for _, child := range block.Children {
			r.WriteString(child.zhubaiBlockMarkdown())
		}
		r.WriteString("\n")
	case "bulleted-list":
		for _, child := range block.Children {
			r.WriteString(fmt.Sprintf("  - %s\n", child.zhubaiBlockMarkdown()))
		}
		r.WriteString("\n")
	case "numbered-list":
		for i, child := range block.Children {
			r.WriteString(fmt.Sprintf("  %d. %s\n", i, child.zhubaiBlockMarkdown()))
		}
		r.WriteString("\n")
	case "block-quote":
		r.WriteString(fmt.Sprintf("> %s", zhubaiBlockListMarkdown(block.Children)))
		r.WriteString("\n")
	case "block-code":
		r.WriteString(fmt.Sprintf("```\n%s\n```", zhubaiBlockListMarkdown(block.Children)))
		r.WriteString("\n")
	case "divider":
		r.WriteString("---\n")
	case "heading-one":
		r.WriteString(fmt.Sprintf("# %s ", block.Children[0].Text))
		r.WriteString("\n")
	case "heading-two":
		r.WriteString(fmt.Sprintf("## %s ", block.Children[0].Text))
		r.WriteString("\n")
	case "list-item":
		return zhubaiBlockListMarkdown(block.Children)
	case "link":
		r.WriteString(fmt.Sprintf("[%s](%s)", zhubaiBlockListMarkdown(block.Children), block.URL))
		r.WriteString("\n")
	case "image":
		r.WriteString(fmt.Sprintf("![%s](%s)", zhubaiBlockListMarkdown(block.Children), block.URL))
		r.WriteString("\n")
	default:
		if block.Bold && block.Italic {
			r.WriteString(fmt.Sprintf(" **_%s_** ", block.Text))
		} else if block.Bold {
			r.WriteString(fmt.Sprintf(" **%s** ", block.Text))
		} else if block.Italic {
			r.WriteString(fmt.Sprintf(" _%s_ ", block.Text))
		} else {
			r.WriteString(block.Text)
		}
	}

	return r.String()
}
