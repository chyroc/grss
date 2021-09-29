package helper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// <a href="http://xueqiu.com/S/SZ000799" target="_blank">$酒鬼酒(SZ000799)$</a>
// 视频B站同步发了，欢迎关注。是的，我也有B站账号的 <span class="url-icon"><img alt=[允悲] src="xxx" /></span> <a data-url=\"http://t.cn/A6MtWp6T\" href=\"http://t.cn/A6MtWp6T\" data-hide=\"\"><span class='url-icon'><img style='width: 1rem;height: 1rem' src='https://h5.sinaimg.cn/upload/2015/09/25/3/timeline_card_small_web_default.png'></span><span class=\"surl-text\">只要5分钟，在云服务器上一键安装180个优质云应用</span></a> ",
func ToTitleText(s string, size int, padding string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err == nil {
		s = doc.Text()
	}
	ss := []rune(s) // not []byte
	if len(ss) <= size {
		return s
	}
	return string(ss[:size]) + padding
}

func ToHtml(s string) string {
	res := []string{}
	for _, v := range strings.Split(s, "\n") {
		res = append(res, fmt.Sprintf("<div>%s</div>", v))
	}
	return fmt.Sprintf("<div>\n%s\n</div>", strings.Join(res, "\n"))
}
