package grss

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chyroc/grss/interface/routers"
)

func updateReadme() error {
	items := []string{}
	for _, v := range routers.Get() {
		items = append(items, fmt.Sprintf("- [%s](/xml%s.xml): %s", v.Path, v.Path, v.Source.Title))
	}

	data := fmt.Sprintf(`# Git RSS

## Routers

%s
`, strings.Join(items, "\n"))
	return ioutil.WriteFile("./index.md", []byte(data), 0o666)
}
