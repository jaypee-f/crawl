package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// TODO
// Add concurrency will need it later
// Split into packages and testing

func main() {
	fmt.Println("let's crawl")
	if len(os.Args) < 2 {
		fmt.Println("provide a url to crawl")
	}

	linksQueue := make(chan string)
	seenLinks := make(chan string)

	go func() {
		linksQueue <- os.Args[1]
	}()
	go filterLinks(linksQueue, seenLinks)

	for url := range seenLinks {
		ConsumeLinksQueue(url, linksQueue)
	}

}

func filterLinks(allLinks chan string, unSeenLinks chan string) {
	var crawled = make(map[string]time.Time)
	for link := range allLinks {
		_, ok := crawled[link]
		if !ok {
			crawled[link] = time.Now()
			unSeenLinks <- link
		} else {
			fmt.Println("SEEN", link)
		}
	}
}

func ConsumeLinksQueue(url string, queue chan string) {
	links := GetLinksFromUrl(url)
	for i := range links {
		if links[i] != "" {
			link := filterLink(links[i], url)
			//fmt.Println(link)
			go func() { queue <- link }()
		}
	}
}

func GetLinksFromUrl(url string) []string {
	r, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return getLinks(r.Body)
}

func getLinks(body io.ReadCloser) []string {
	defer func() {
		body.Close()
	}()

	t := html.NewTokenizer(body)
	links := make([]string, 0)

	for {
		switch t.Next() {
		case html.ErrorToken:
			if len(links) > 0 {
				return links
			}
			return nil
		case html.StartTagToken:
			link := getLink(t.Token())
			if link != "" {
				links = append(links, link)
			}
		}
	}
}

func getLink(t html.Token) string {
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func filterLink(link, page string) string {
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
	return pageUrl.ResolveReference(uri).String()

}
