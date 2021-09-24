package ycombinator_best

import (
	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/routers/ycombinator/internal"
)

func New(map[string]string) (*fetch.Source, error) {
	return ycombinator_internal.Generate("Hacker News - Top Links", "https://news.ycombinator.com/best")(nil)
}
