package grss

import (
	"flag"
	"log"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	"github.com/chyroc/grss/internal/routers"
)

func Run() {
	path := ""
	flag.StringVar(&path, "path", "", "path")
	flag.Parse()

	routers := routers.Get()
	log.Printf("load %d router", len(routers))

	removeOldData() // 删除老数据

	for _, router := range routers {
		if path != "" && router.Path != path {
			log.Printf("[grss] skip %q", router.Path)
			continue
		} else {
			log.Printf("[grss] gen %q", router.Path)
		}
		sourceGetter := router.Source

		if len(router.Args) == 0 {
			router.Args = []map[string]string{nil}
		}
		for _, args := range router.Args {
			path := helper.ToJoinArgsURL(router.Path, args)

			feed, err := fetch.Fetch(sourceGetter, args)
			if err != nil {
				log.Printf("fetch %s failed: %s", path, err)
				continue
			}

			if err := DumpFeed(path, feed); err != nil {
				log.Printf("dump feed %s failed: %s", feed.Title, err)
				continue
			}
		}
	}

	if err := updateReadme(); err != nil {
		log.Printf("update readme failed: %s", err)
	}
}
