package grss

import (
	"log"

	"github.com/chyroc/grss/interface/fetch"
	"github.com/chyroc/grss/interface/routers"
)

func Run() {
	routers, err := loadRouters()
	if err != nil {
		panic(err)
	}
	log.Printf("load %d router", len(routers))

	for _, router := range routers {
		path := router.Path
		source := router.Source
		log.Printf("start fetch %s at %s", source.Title, path)
		feed, err := fetch.Fetch(source)
		if err != nil {
			log.Printf("fetch %s failed: %s", source.Title, err)
			continue
		}

		if err := DumpFeed(path, feed); err != nil {
			log.Printf("dump feed %s failed: %s", source.Title, err)
			continue
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
