package helper

import (
	"os"
	"strings"

	"github.com/chyroc/baidufanyi"
	"github.com/chyroc/go-lambda"
)

func Fanyi(text string) (string, error) {
	translateResult, err := baidufanyi.New(baidufanyi.WithCredential(os.Getenv("BAIDUFANYI_APP_ID"), os.Getenv("BAIDUFANYI_APP_SECRET"))).Translate(text, baidufanyi.LanguageEn, baidufanyi.LanguageZh)
	if err != nil {
		return "", err
	}
	translateResultString, _ := lambda.New(translateResult).MapList(func(idx int, obj interface{}) interface{} {
		return obj.(*baidufanyi.TranslateResult).Dst
	}).ToStringSlice()

	return strings.Join(translateResultString, "\n"), nil
}

func FanyiAndAppend(text, sep string) string {
	fanyi, err := Fanyi(text)
	if err != nil {
		return text
	}

	return strings.Join([]string{fanyi, text}, sep)
}
