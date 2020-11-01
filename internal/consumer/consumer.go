package consumer

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jaypee-f/crawler/internal/fetching"
)

type Crawler struct {
	base string
}

func New(base string) *Crawler {
	return &Crawler{base: base}
}

func (c *Crawler) Start() {
	linksQueue := make(chan string)
	seenLinks := make(chan string)
	go func() {
		linksQueue <- c.base
	}()
	go c.filterLinks(linksQueue, seenLinks)

	for url := range seenLinks {
		c.consumeLinksQueue(url, linksQueue)
	}
}

func (c *Crawler) filterLinks(allLinks chan string, unSeenLinks chan string) {
	var crawled = make(map[string]time.Time)
	for link := range allLinks {
		_, ok := crawled[link]
		if !ok {
			crawled[link] = time.Now()
			fmt.Println(link)
			unSeenLinks <- link
		}
	}
}

func (c *Crawler) consumeLinksQueue(url string, queue chan string) {
	links := fetching.GetLinksFromUrl(url)
	for i := range links {
		link := c.filterLink(links[i], url)
		if link != "" {
			go func() { queue <- link }()
		}
	}
}

func (c *Crawler) filterLink(link, page string) string {
	// Is it one of ours?

	// remove hash
	l := strings.Split(link, "#")

	uri, err := url.Parse(l[0])
	if err != nil {
		return ""
	}
	pageUrl, err := url.Parse(page)
	if err != nil {
		return ""
	}
	link = pageUrl.ResolveReference(uri).String()

	if !strings.HasPrefix(link, c.base) {
		return ""
	}

	return link
}
