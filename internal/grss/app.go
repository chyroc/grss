package grss

import (
	"log"

	"github.com/chyroc/grss/internal/fetch"
	"github.com/chyroc/grss/internal/helper"
	"github.com/chyroc/grss/internal/routers"
)

func Run() {
	routers, err := loadRouters()
	if err != nil {
		panic(err)
	}
	log.Printf("load %d router", len(routers))

	for _, router := range routers {
		sourceGetter := router.Source

		if len(router.Args) == 0 {
			router.Args = []map[string]string{nil}
		}
		for _, args := range router.Args {
			path := helper.JoinArgsURL(router.Path, args)

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

func loadRouters() ([]routers.Router, error) {
	routers := routers.Get()

	return routers, nil
}
