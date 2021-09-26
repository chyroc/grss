package helper

import (
	"crypto/md5"
	"fmt"
)

func Md5(s string) string {
	ins := md5.New()
	ins.Write([]byte(s))
	return fmt.Sprintf("%x", ins.Sum(nil))
}
