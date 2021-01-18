package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	link "github.com/SirNoob97/gophercises/html-link-parser"
)

const xmlns = "https://www.sitemaps.org/schemas/sitemap/9.0"

type loc struct {
	Loc string `xml:"loc"`
}

type urlSet struct {
	URLs  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximun number of links deep to traverse")
	flag.Parse()

	views := bfs(*urlFlag, *maxDepth)

	toXML := urlSet{
		Xmlns: xmlns,
	}
	for _, view := range views {
		toXML.URLs = append(toXML.URLs, loc{view})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")

	if err := enc.Encode(toXML); err != nil {
		exit(err.Error())
	}

	fmt.Println()
}

// check if you have already visited the url
func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = struct{}{}

			for _, link := range getLinks(url) {
					nq[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func getLinks(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		exit(err.Error())
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}

	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func filter(links []string, keepFunc func(string) bool) []string {
	var ret []string
	for _, l := range links {
		if keepFunc(l) {
			ret = append(ret, l)
		}
	}
	return ret
}

func hrefs(r io.Reader, base string) []string {
	links, err := link.Parser(r)
	if err != nil {
		return []string{}
	}

	var ret []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func exit(msg string) {
	log.Fatalf(msg)
}
