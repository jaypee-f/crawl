package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawler_Start(t *testing.T) {
	//	c := Crawler{}
	//c.Start()
}

func TestCrawler_dedupeLinks(t *testing.T) {
	c := New("")
	c.links <- "foo"
	c.links <- "bar"
	c.links <- "foo"
	close(c.links)

	c.dedupeLinks(nil)

	links := make([]string, 0)
	for {
		link, ok := <-c.unSeenLinks
		if !ok {
			break
		}
		links = append(links, link)
	}
	assert.Equal(t, []string{"foo", "bar"}, links)
}

func TestCrawler_sanitiseLink(t *testing.T) {

	type test struct {
		name   string
		base   string
		link   string
		page   string
		result string
	}
	tests := []test{
		{
			name:   "Our page - sucess",
			base:   "http://example.com",
			link:   "http://example.com",
			page:   "http://example.com",
			result: "http://example.com",
		},
		{
			name:   "Our page - relative link",
			base:   "http://example.com",
			link:   "/foo",
			page:   "http://example.com",
			result: "http://example.com/foo",
		},
		{
			name:   "Our page - relative link hash check",
			base:   "http://example.com",
			link:   "/foo#bar",
			page:   "http://example.com",
			result: "http://example.com/foo",
		},
		{
			name:   "Not out page - relative link hash check",
			base:   "http://example.com",
			link:   "/foo#bar",
			page:   "http://bar.com",
			result: "",
		},
		{
			name:   "bad link",
			base:   "",
			link:   "http://x:namedport",
			page:   "",
			result: "",
		},
		{
			name:   "bad base",
			base:   "http://x:namedport",
			link:   "http://example.com",
			page:   "http://example.com",
			result: "",
		},
		{
			name:   "bad base",
			base:   "http://example.com",
			link:   "http://example.com",
			page:   "http://x:namedport",
			result: "",
		},
		{
			name:   "nothing",
			base:   "",
			link:   "",
			page:   "",
			result: "",
		},
	}

	for _, test := range tests {
		c := Crawler{base: test.base}
		assert.Equal(t, test.result, c.santiseLink(test.link, test.page), test.name)
	}
}
