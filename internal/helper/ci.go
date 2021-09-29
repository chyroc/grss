package helper

import (
	"os"
)

var IsInCI bool

func init() {
	IsInCI = os.Getenv("IN_CI") != ""
}
