package pingwest_status_test

import (
	"testing"

	"github.com/chyroc/go-assert"
	"github.com/chyroc/grss/interface/fetch"
	"github.com/chyroc/grss/interface/routers/pingwest/status"
	"github.com/davecgh/go-spew/spew"
)

func TestName(t *testing.T) {
	as := assert.New(t)
	as.Nil(nil)

	feed, err := fetch.Fetch(pingwest_status.New())
	as.Nil(err)
	spew.Dump(feed)
}
