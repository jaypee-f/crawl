package fetching

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

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
