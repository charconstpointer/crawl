package main

import (
	"flag"
	"fmt"

	"github.com/charconstpointer/crawl/pkg/crawler"
	"go.uber.org/zap"
)

func main() {
	root := flag.String("root", "https://google.com", "foo")
	phrase := flag.String("phrase", "foo", "phrase to find")
	rl := flag.Int("rl", 5, "rate limit")
	logger, _ := zap.NewProduction()

	flag.Parse()
	logger.Sugar().Infow("Flags %s, %s, %d", zap.String("root", *root), zap.String("phrase", *phrase), zap.Int("rl", *rl))
	c := crawler.NewCrawler(*rl)
	go func(c *crawler.Crawler) {
		for {
			select {
			case msg := <-c.C:
				fmt.Printf(">%s\n", msg)
			}
		}
	}(c)
	err := c.Crawl(*root, *phrase)
	if err != nil {
		logger.Sugar().Error(err.Error())
		return
	}
}
