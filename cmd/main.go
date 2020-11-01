package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("let's crawl")
	if len(os.Args) < 2 {
		fmt.Println("provide a url to crawl")
	}


	links := GetLinksFromUrl(os.Args[1])
	if links != nil {
		fmt.Println(links)
	}
}

func getLinks(body io.ReadCloser) []string{
	t := html.NewTokenizer(body)
	links := make([]string,0)

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
				links = append(links,link)
			}
		}
	}


}

func GetLinksFromUrl(url string)  []string {
	r, err := http.Get(url)
	defer func () {
		if r != nil {
			r.Body.Close()
		}
	}()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return getLinks(r.Body)
}

func getLink(t html.Token) string {
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}