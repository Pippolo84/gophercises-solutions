package linkextract

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link is a struct holding the info extracted
// from the anchors in the parsed HTML
type Link struct {
	Href string
	Text string
}

// Links takes a pointer to a bufio.Reader where it
// expects to find an UTF-8 encoded HTML document.
// It then extracts all the links from that documents
// and returns them.
// If an error occurs, an empty Link slices along with
// the error itself will be returned.
func Links(r io.Reader) ([]Link, error) {
	var links []Link

	htmlDoc, err := html.Parse(r)
	if err != nil {
		return links, err
	}

	visitAnchor(htmlDoc, &links)

	return links, nil
}

func visitAnchor(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var link Link

		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link.Href = attr.Val
			}
		}

		var textBuilder strings.Builder
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visitText(c, &textBuilder)
		}

		link.Text = strings.TrimSpace(textBuilder.String())

		*links = append(*links, link)
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visitAnchor(c, links)
		}
	}
}

func visitText(n *html.Node, textBuilder *strings.Builder) {
	if n.Type == html.TextNode {
		textBuilder.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visitText(c, textBuilder)
	}
}
