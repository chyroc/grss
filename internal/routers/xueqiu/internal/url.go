package internal

import (
	"fmt"
	"strings"
)

func JoinURL(target string) string {
	if strings.HasPrefix(target, "https://xueqiu.com") || strings.HasPrefix(target, "http://xueqiu.com") {
		return target
	}
	return fmt.Sprintf("https://xueqiu.com%s", target)
}
