package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

var (
	root = flag.String("root", "https://wykop.pl", "root")
)

func main() {
	flag.Parse()
	log.Println(*root)
	c := http.Client{}
	crawl(c, *root)

}

func crawl(c http.Client, root string) {
	r, err := http.NewRequest("GET", root, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	res, err := c.Do(r)
	if err != nil {
		log.Fatal(err.Error())
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	exp, _ := regexp.Compile("https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)")
	matches := exp.FindAllString(string(content), -1)
	for _, m := range matches {
		log.Printf("calling %s found %d matches", m, len(matches))
		crawl(c, m)
		time.Sleep(5 * time.Second)
	}
}
