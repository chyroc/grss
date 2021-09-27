package helper

import (
	"github.com/chyroc/gorequests"
)

var Req *gorequests.Factory

func init() {
	Req = gorequests.NewFactory(
		gorequests.WithLogger(gorequests.NewDiscardLogger()),
		gorequests.WithHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"),
	)
}
