package crawler

import (
	"gophercises/link/linkextract"
	"gophercises/sitemap/bfs"
	"net/http"
	"strings"
)

// Crawl return all links in the webpage specified by domain
// Only the link pointing to the same domain will be returned
func Crawl(domain string) ([]string, error) {
	startHref := href{
		domain: domain,
		path:   "/",
	}

	generators, err := bfs.Bfs(startHref, bfs.VisitorFunc(hrefVisitor))
	if err != nil {
		return nil, err
	}

	var links []string
	for _, g := range generators {
		links = append(links, strings.TrimPrefix(g.ID(), domain))
	}

	return links, err
}

type href struct {
	domain string
	path   string
}

func (href href) ID() string {
	if strings.HasPrefix(href.path, "/") {
		return href.domain + href.path
	} else {
		return href.domain + "/" + href.path
	}
}

func (href href) Value() interface{} {
	return href.ID()
}

func hrefVisitor(node bfs.Producer) ([]bfs.Producer, error) {
	var next []bfs.Producer
	url := node.ID()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	links, err := linkextract.Links(resp.Body)
	if err != nil {
		return next, nil
	}

	for _, link := range links {
		// check if link goes inside the same domain
		if strings.HasPrefix(link.Href, "http") && !strings.HasPrefix(link.Href, node.(href).domain) {
			continue
		}

		var path string
		if strings.HasPrefix(link.Href, node.(href).domain) {
			// remove domain if present
			path = strings.TrimPrefix(link.Href, node.(href).domain)
		} else {
			path = link.Href
		}

		next = append(next, href{
			domain: node.(href).domain,
			path:   path,
		})
	}

	return next, nil
}