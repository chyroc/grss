package helper

import (
	"github.com/chyroc/gorequests"
)

var Req *gorequests.Factory

func init() {
	Req = gorequests.NewFactory(gorequests.WithLogger(gorequests.NewDiscardLogger()))
}
