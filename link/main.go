package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Link is a struct holding the info extracted
// from the anchors in the parsed HTML
type Link struct {
	Href string
	Text string
}

func main() {
	var fileName string

	flag.StringVar(&fileName, "file", "input.html", "HTML file to parse")
	flag.Parse()

	in, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(in)

	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var links []Link
	visitAnchor(doc, &links)

	for _, link := range links {
		fmt.Println(link)
	}
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
