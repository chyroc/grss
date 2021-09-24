package ycombinator_newest

import (
	"github.com/chyroc/grss/internal/fetch"
	ycombinator_internal "github.com/chyroc/grss/internal/routers/ycombinator/internal"
)

func New(map[string]string) (*fetch.Source, error) {
	return ycombinator_internal.Generate("Hacker News - New Links", "https://news.ycombinator.com/newest")(nil)
}
