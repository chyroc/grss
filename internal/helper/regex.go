package helper

import (
	"regexp"
)

func GetOneMatchString(s string, regexp *regexp.Regexp) string {
	m := regexp.FindStringSubmatch(s)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}
