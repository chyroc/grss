package helper

import (
	"strings"
)

func ToJoinArgsURL(url string, args map[string]string) string {
	for k, v := range args {
		url = strings.ReplaceAll(url, ":"+k, v)
	}
	return url
}
