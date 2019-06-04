package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"gophercises/link/linkextract"
	"log"
	"net/http"
	"strings"
)

const (
	SitemapXmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

// FIXME: use a different package for those urlset, url, etc.

type SitemapUrl struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

type SitemapUrlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []SitemapUrl
}

func main() {
	var url string

	flag.StringVar(&url, "url", "http://www.ansa.it", "URL of the website to map")
	flag.Parse()

	paths, err := crawlDomain(url)
	if err != nil {
		log.Fatal(err)
	}

	smUrlset := SitemapUrlset{Xmlns: SitemapXmlns}
	for _, path := range paths {
		smURL := SitemapUrl{Loc: url + path}
		smUrlset.Urls = append(smUrlset.Urls, smURL)
	}

	xmlString, err := xml.MarshalIndent(smUrlset, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	xmlDoc := append([]byte(xml.Header), xmlString...)

	fmt.Println(string(xmlDoc))
}

func crawlDomain(domain string) ([]string, error) {
	visited := make(map[string]bool)
	var hrefs []string
	var queue []string

	queue = append(queue, "/")
	for len(queue) > 0 {
		var path string

		path, queue = queue[0], queue[1:]
		if _, ok := visited[path]; ok {
			continue
		}

		resp, err := http.Get(domain + path)
		if err != nil {
			return hrefs, err
		}
		defer resp.Body.Close()

		visited[path] = true

		hrefs = append(hrefs, path)

		links, err := linkextract.Links(resp.Body)

		for _, link := range links {
			// check if link goes to the same domain
			if strings.HasPrefix(link.Href, "http") && !strings.HasPrefix(link.Href, domain) {
				continue
			}

			// remove domain if present
			if strings.HasPrefix(link.Href, domain) {
				path = strings.TrimPrefix(link.Href, domain)
			}

			queue = append(queue, path)
		}
	}

	return hrefs, nil
}
