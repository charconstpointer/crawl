package main

import (
	"flag"
	"fmt"

	"github.com/charconstpointer/crawl/pkg/crawler"
	"github.com/labstack/gommon/log"
)

func main() {
	root := flag.String("root", "https://wykop.pl", "romek898")
	phrase := flag.String("phrase", "foo", "phrase to find")
	rl := flag.Int("rl", 5, "rate limit")

	flag.Parse()
	log.Infof("Flags %s, %s, %d", *root, *phrase, *rl)
	c := crawler.NewCrawler(*rl)
	go func(c *crawler.Crawler) {
		for {
			select {
			case msg := <-c.C:
				fmt.Printf(">%s\n", msg)
			}
		}
	}(c)
	_, err := c.Crawl(*root, *phrase)
	if err != nil {
		log.Error(err)
		return
	}
}
