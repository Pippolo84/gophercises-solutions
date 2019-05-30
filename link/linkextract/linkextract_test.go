package linkextract

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestLinksSingleLink(t *testing.T) {
	document := `
<html>
<body>
	<h1>Hello!</h1>
	<a href="/other-page">A link to another page</a>
</body>
</html>`

	links, err := Links(strings.NewReader(document))
	if err != nil {
		t.Errorf("Call to Links failed with error %v\n", err)
	}

	expected := []Link{
		Link{
			Href: "/other-page",
			Text: "A link to another page",
		},
	}

	if len(links) != len(expected) {
		t.Errorf("Expected %v links, got %v\n", len(expected), len(links))
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected %v link to be %v, got %v\n", i, expected[i], link)
		}
	}
}

func TestLinksMultipleLinks(t *testing.T) {
	document := `
<html>
<head>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
	<h1>Social stuffs</h1>
	<div>
	<a href="https://www.twitter.com/joncalhoun">
		Check me out on twitter
		<i class="fa fa-twitter" aria-hidden="true"></i>
	</a>
	<a href="https://github.com/gophercises">
		Gophercises is on <strong>Github</strong>!
	</a>
	</div>
</body>
</html>`

	links, err := Links(strings.NewReader(document))
	if err != nil {
		t.Errorf("Call to Links failed with error %v\n", err)
	}

	expected := []Link{
		Link{
			Href: "https://www.twitter.com/joncalhoun",
			Text: "Check me out on twitter",
		},
		Link{
			Href: "https://github.com/gophercises",
			Text: "Gophercises is on Github!",
		},
	}

	if len(links) != len(expected) {
		t.Errorf("Expected %v links, got %v\n", len(expected), len(links))
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected %v link to be %v, got %v\n", i, expected[i], link)
		}
	}
}

func TestLinksLinkWithComment(t *testing.T) {
	document := `
<html>
<body>
	<a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>`

	links, err := Links(strings.NewReader(document))
	if err != nil {
		t.Errorf("Call to Links failed with error %v\n", err)
	}

	expected := []Link{
		Link{
			Href: "/dog-cat",
			Text: "dog cat",
		},
	}

	if len(links) != len(expected) {
		t.Errorf("Expected %v links, got %v\n", len(expected), len(links))
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected %v link to be %v, got %v\n", i, expected[i], link)
		}
	}
}

func TestVisitAnchorEmpty(t *testing.T) {
	node := html.Node{
		Type: html.DocumentNode,
	}

	var links []Link
	visitAnchor(&node, &links)

	if len(links) > 0 {
		t.Errorf("Expected empty links slice, got %v\n", links)
	}
}

func TestVisitAnchorSingle(t *testing.T) {
	childNode := html.Node{
		Type: html.TextNode,
		Data: "test text",
	}
	parentNode := html.Node{
		FirstChild: &childNode,
		LastChild:  &childNode,
		Type:       html.ElementNode,
		Data:       "a",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "href",
				Val: "/test",
			},
		},
	}

	var links []Link
	visitAnchor(&parentNode, &links)

	if len(links) != 1 {
		t.Errorf("Expected one link in slice, got %v\n", len(links))
	}

	expected := Link{
		Href: "/test",
		Text: "test text",
	}
	if links[0] != expected {
		t.Errorf("Expected link %v, got %v\n", expected, links[0])
	}
}

func TestVisitAnchorTrimText(t *testing.T) {
	childNode := html.Node{
		Type: html.TextNode,
		Data: "\n    test text\t \t\n    ",
	}
	parentNode := html.Node{
		FirstChild: &childNode,
		LastChild:  &childNode,
		Type:       html.ElementNode,
		Data:       "a",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "href",
				Val: "/test",
			},
		},
	}

	var links []Link
	visitAnchor(&parentNode, &links)

	if len(links) != 1 {
		t.Errorf("Expected one link in slice, got %v\n", len(links))
	}

	expected := Link{
		Href: "/test",
		Text: "test text",
	}
	if links[0] != expected {
		t.Errorf("Expected link %v, got %v\n", expected, links[0])
	}
}
func TestVisitAnchorSiblings(t *testing.T) {
	secondChildNode := html.Node{
		Type: html.TextNode,
		Data: "second test text",
	}
	firstChildNode := html.Node{
		NextSibling: &secondChildNode,

		Type: html.TextNode,
		Data: "first test text",
	}
	parentNode := html.Node{
		FirstChild: &firstChildNode,
		LastChild:  &secondChildNode,

		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "href",
				Val: "/test",
			},
		},
	}

	var links []Link
	visitAnchor(&parentNode, &links)

	if len(links) != 1 {
		t.Errorf("Expected one link in slice, got %v\n", len(links))
	}

	expected := Link{
		Href: "/test",
		Text: "first test textsecond test text",
	}

	if links[0] != expected {
		t.Errorf("Expected link %v, got %v\n", expected, links[0])
	}
}

func TestVisitAnchorMultiple(t *testing.T) {
	secondTextNode := html.Node{
		Type: html.TextNode,
		Data: "second test text",
	}
	secondElementNode := html.Node{
		FirstChild: &secondTextNode,
		LastChild:  &secondTextNode,

		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "href",
				Val: "/test2",
			},
		},
	}

	firstTextNode := html.Node{
		Type: html.TextNode,
		Data: "first test text",
	}
	firstElementNode := html.Node{
		NextSibling: &secondElementNode,
		FirstChild:  &firstTextNode,
		LastChild:   &firstTextNode,

		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "href",
				Val: "/test1",
			},
		},
	}

	documentNode := html.Node{
		FirstChild: &firstElementNode,
		LastChild:  &secondElementNode,

		Type: html.DocumentNode,
	}

	var links []Link
	visitAnchor(&documentNode, &links)

	expected := []Link{
		Link{
			Href: "/test1",
			Text: "first test text",
		},
		Link{
			Href: "/test2",
			Text: "second test text",
		},
	}

	if len(links) != len(expected) {
		t.Errorf("Expected %v link in slice, got %v\n", len(expected), len(links))
	}

	for i := 0; i < len(expected); i++ {
		if links[i] != expected[i] {
			t.Errorf("Expected %v-th link %v, got %v\n", i, expected[i], links[i])
		}
	}
}

func TestVisitText(t *testing.T) {
	expected := "test text"

	node := html.Node{
		Type: html.TextNode,
		Data: expected,
	}

	var textBuilder strings.Builder
	visitText(&node, &textBuilder)

	result := textBuilder.String()
	if result != expected {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestVisitTextEmpty(t *testing.T) {
	expected := ""

	node := html.Node{
		Type: html.TextNode,
		Data: expected,
	}

	var textBuilder strings.Builder
	visitText(&node, &textBuilder)

	result := textBuilder.String()
	if result != expected {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}
