package grss

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chyroc/grss/internal/helper"
	"github.com/chyroc/grss/internal/routers"
)

func updateReadme() error {
	items := []string{}
	for _, v := range routers.Get() {
		if len(v.Args) == 0 {
			v.Args = []map[string]string{nil}
		}
		for _, args := range v.Args {
			source, err := v.Source(args)
			if err != nil {
				panic(err)
			}
			path := helper.ToJoinArgsURL(v.Path, args)
			items = append(items, fmt.Sprintf("| %s | %s | %s | %s |",
				source.Title,
				fmt.Sprintf("[JSON](./json/latest%s.json)", path),
				fmt.Sprintf("[RSS](./xml%s.xml)", path),
				fmt.Sprintf("[HTML](./html%s/index.html)", path),
			))
		}
	}

	data := fmt.Sprintf(`![](./header.png)

# Git RSS

Homepage: https://rss.chyroc.cn

Generate Tool: https://github.com/chyroc/grss

## Routers

| Title | JSON | RSS | HTML |
| ---   | ---  | --- | ---  |
%s
`, strings.Join(items, "\n"))

	if err := ioutil.WriteFile("./README.md", []byte(data), 0o666); err != nil {
		panic(err)
	}

	return nil
}
