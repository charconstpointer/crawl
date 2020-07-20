package crawler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"go.uber.org/ratelimit"
)

type crawler interface {
	Crawl(URL string, phrase string) error
}

type Crawler struct {
	rl      ratelimit.Limiter
	C       chan string
	visited []string
	mutex   sync.Mutex
}

func NewCrawler(n int) *Crawler {
	return &Crawler{ratelimit.New(n), make(chan string), make([]string, 0), sync.Mutex{}}
}

func (c *Crawler) Crawl(URL string, phrase string) error {
	err := c.walk(URL, phrase, 0)
	return err
}

func (c *Crawler) walk(URL string, phrase string, depth int) error {
	if depth > 3 {
		return errors.New("too deep")
	}
	res, err := http.Get(URL)
	if err != nil {
		return err
	}

	c.addVisited(URL)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	found := search(phrase, string(b), URL)

	if found {
		c.C <- URL
		return nil
	}

	urls := c.findUrls(string(b))
	for _, u := range urls {
		_ = c.rl.Take()
		go func(u string, ch chan<- string) {
			err := c.walk(u, phrase, depth+1)
			if err != nil {

				// fmt.Printf("%v", err)
			}
		}(u, c.C)

	}
	return errors.New("could not find requested phrase")
}

func (c *Crawler) findUrls(source string) []string {
	exp := "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	r, _ := regexp.Compile(exp)
	links := r.FindAllString(source, -1)
	notVisited := make([]string, 0)
	for _, l := range links {
		if !c.hasVisited(l) {
			notVisited = append(notVisited, l)
		}
	}
	return notVisited
}

func search(phrase string, content string, URL string) bool {
	found := strings.Contains(content, phrase)
	return found
}

func (c *Crawler) hasVisited(URL string) bool {
	for _, v := range c.visited {
		if v == URL {
			return true
		}
	}
	return false
}

func (c *Crawler) addVisited(URL string) {
	c.mutex.Lock()
	if !c.hasVisited(URL) {
		c.visited = append(c.visited, URL)
	}
	c.mutex.Unlock()
}
