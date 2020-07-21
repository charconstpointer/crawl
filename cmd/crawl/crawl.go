package main

import (
	"flag"

	"github.com/charconstpointer/crawl/pkg/crawler"
	"github.com/labstack/gommon/log"
)

func main() {
	root := flag.String("root", "https://google.com", "foo")
	phrase := flag.String("phrase", "foo", "phrase to find")
	limit := flag.Int("limit", 2, "depth limit")
	rl := flag.Int("rl", 5, "rate limit")

	flag.Parse()
	log.Infof("flags %s, %s, %d\n", "root", *root, "phrase", *phrase, "rl", *rl)
	c := crawler.NewCrawler(*rl, *limit)
	err := c.Crawl(*root, *phrase)

	if err != nil {
		log.Error(err)
		return
	}

	go func() {
		for {
			select {
			case msg := <-c.C:
				log.Infof(">%s\n", msg)
			}
		}
	}()

}
