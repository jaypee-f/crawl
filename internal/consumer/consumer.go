package consumer

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jaypee-f/crawler/internal/fetching"
)

const maxWait = time.Second * 5

type Crawler struct {
	base        string
	links       chan string
	unSeenLinks chan string
}

func New(base string) *Crawler {
	return &Crawler{
		base:        base,
		links:       make(chan string, 5),
		unSeenLinks: make(chan string, 5),
	}
}

func (c *Crawler) Start() {
	done := make(chan bool)
	go func() {
		c.links <- c.base
	}()

	go c.dedupeLinks(done)

	for i := 0; i < 3; i++ {
		go c.crawl()
	}
	<-done
}

func (c *Crawler) crawl() {
	for {
		select {
		case link, ok := <-c.unSeenLinks:
			if !ok {
				fmt.Println("consumer quits")
				return
			}
			c.consumeLinks(link)
		default:
		}
	}
}

func (c *Crawler) dedupeLinks(done chan bool) {
	crawled := make(map[string]time.Time)

	timer := time.NewTimer(maxWait)

	for {
		select {
		case link, ok := <-c.links:
			if !ok {
				close(c.unSeenLinks)
				return
			}
			timer.Reset(maxWait)
			_, ok = crawled[link]
			if !ok {
				crawled[link] = time.Now()
				c.unSeenLinks <- link
			}
		case <-timer.C:
			fmt.Println("Done")
			done <- true
		}
	}
}

func (c *Crawler) consumeLinks(url string) {
	links := fetching.GetLinksFromUrl(url)
	pageMap := make(map[string]bool)
	for i := range links {
		link := c.santiseLink(links[i], url)
		if link != "" {
			pageMap[link] = true
			go func() { c.links <- link }()
		}
	}

	displayPageLinks(url, pageMap)
}

func displayPageLinks(url string, links map[string]bool) {
	page := make([]string, 0)
	page = append(page, url)

	for link := range links {
		page = append(page, fmt.Sprint("\t", link))
	}
	fmt.Print(strings.Join(page, "\n"), "\n")
}

func (c *Crawler) santiseLink(link, page string) string {
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

	// Is it one of ours?
	if !strings.HasPrefix(link, c.base) {
		return ""
	}

	return link
}
