package crawler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"go.uber.org/ratelimit"
)

type Crawler struct {
	rl ratelimit.Limiter
	C  chan string
}

func NewCrawler(n int) *Crawler {
	return &Crawler{ratelimit.New(n), make(chan string)}
}

func (c *Crawler) Crawl(URL string, phrase string) (string, error) {
	res, err := c.walk(URL, phrase, 0)
	return res, err
}

func (c *Crawler) walk(URL string, phrase string, depth int) (string, error) {
	if depth > 3 {
		return "", errors.New("too deep")
	}
	res, err := http.Get(URL)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	found := search(phrase, string(b), URL)

	if found {
		return URL, nil
	}

	urls := findUrls(string(b))
	for _, u := range urls {
		_ = c.rl.Take()
		go func(u string, ch chan<- string) {
			res, _ := c.walk(u, phrase, depth+1)
			// fmt.Printf("%d ] %s ~> %s\n", depth, u, res)
			if res != "" {
				ch <- res
			}
		}(u, c.C)

	}
	return "", errors.New("could not find requested phrase")
}

func findUrls(source string) []string {
	exp := "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	r, _ := regexp.Compile(exp)
	links := r.FindAllString(source, -1)
	return links
}

func search(phrase string, content string, URL string) bool {
	found := strings.Contains(content, phrase)
	return found
}
