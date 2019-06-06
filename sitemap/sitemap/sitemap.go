package sitemap

import (
	"encoding/xml"
)

const (
	sitemapXmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type sitemapURL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

type sitemapUrlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []sitemapURL
}

// Sitemap builds the xml representation of the website
// and returns it as a slice of bytes.
func Sitemap(url string, paths []string) ([]byte, error) {
	smUrlset := sitemapUrlset{Xmlns: sitemapXmlns}
	for _, path := range paths {
		smURL := sitemapURL{Loc: url + path}
		smUrlset.Urls = append(smUrlset.Urls, smURL)
	}

	xmlDoc, err := xml.MarshalIndent(smUrlset, "", "    ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), xmlDoc...), nil
}
