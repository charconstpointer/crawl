package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/charconstpointer/crawl/pkg/set"
	"github.com/labstack/gommon/log"
	"go.uber.org/ratelimit"
)

type crawler interface {
	Crawl(URL string, phrase string, depth int) error
}

type Crawler struct {
	rl      ratelimit.Limiter
	root    string
	limit   int
	C       chan string
	visited set.Set
	mutex   sync.Mutex
}

func NewCrawler(n int, limit int) *Crawler {
	return &Crawler{rl: ratelimit.New(n), C: make(chan string), visited: set.NewSet(), limit: limit}
}

func (c *Crawler) Crawl(URL string, phrase string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	c.root = u.Host
	err = c.walk(URL, phrase, 0)
	return err
}

func (c *Crawler) walk(URL string, phrase string, depth int) error {
	if depth > c.limit {
		return fmt.Errorf("you've reached declared depth limit %d\n", c.limit)
	}

	res, err := http.Get(URL)
	if err != nil {
		log.Error(err)
	}
	c.addVisited(URL)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	urls := c.findUrls(string(b))

	found := search(phrase, string(b))

	if found {
		c.C <- URL
		return nil
	}

	for _, u := range urls {
		_ = c.rl.Take()
		go func(u string, ch chan<- string) {
			_ = c.walk(u, phrase, depth+1)
		}(u, c.C)

	}
	return errors.New("could not find requested phrase")
}

func (c *Crawler) findUrls(source string) []string {
	urlExp := "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	imgExp := "(https?:\\/\\/.*\\.(?:png|jpg|js))"
	r, err := regexp.Compile(urlExp)
	ie, err := regexp.Compile(imgExp)
	if err != nil {
		return nil
	}

	links := r.FindAllString(source, -1)
	notVisited := make([]string, 0)
	for _, l := range links {

		if !strings.Contains(l, c.root) {
			continue
		}

		if ie.Match([]byte(l)) == true {
			continue
		}

		if !c.visited.Contains(l) {
			notVisited = append(notVisited, l)
		}
	}

	return notVisited
}

func search(phrase string, content string) bool {
	found := strings.Contains(content, phrase)
	return found
}

func (c *Crawler) addVisited(URL string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.visited.Contains(URL) {
		c.visited.Add(URL)
	}
}
