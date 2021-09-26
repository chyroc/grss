package helper

import (
	"github.com/PuerkitoBio/goquery"
)

func Selection2List(s *goquery.Selection) (res []*goquery.Selection) {
	s.Each(func(i int, s *goquery.Selection) {
		res = append(res, s)
	})
	return res
}
