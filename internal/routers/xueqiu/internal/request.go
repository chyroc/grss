package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chyroc/grss/internal/helper"
)

func Request(fromLink, apiLink string, resp interface{}) error {
	res, err := helper.Req.New(http.MethodGet, fromLink).Response()
	if err != nil {
		return err
	}
	cookies := res.Header.Values("set-cookie")

	fmt.Println(parseCookie(cookies))
	req := helper.Req.New(http.MethodGet, apiLink).WithHeader("cookie", parseCookie(cookies))
	errResp := struct {
		Err string `json:"error_description"`
	}{}
	text, err := req.Text()
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(text), &errResp); err != nil {
		return err
	} else if errResp.Err != "" {
		return fmt.Errorf(errResp.Err)
	}

	return json.Unmarshal([]byte(text), resp)
}

func parseCookie(cookies []string) (res string) {
	for _, cookie := range cookies {
		v := strings.Split(cookie, ";")
		if len(v) == 0 {
			continue
		}
		res += v[0] + ";"
	}
	return res
}
