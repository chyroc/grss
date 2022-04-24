package helper

import (
	"strings"

	"github.com/chyroc/googletranslate"
)

func Fanyi(text string) (string, error) {
	return googletranslate.Translate(text, googletranslate.En, googletranslate.Zh)
}

func FanyiAndAppend(text, sep string) string {
	fanyi, err := Fanyi(text)
	if err != nil || strings.TrimSpace(fanyi) == "" {
		return text
	}

	return strings.Join([]string{fanyi, text}, sep)
}

func FetchFeedBinAndFanyiAndAppend(url string) string {
	text := AddFeedbinPage(url)
	if text == "" {
		return ""
	}

	return FanyiAndAppend(text, " <br>\n<br>\n原文: ")
}
