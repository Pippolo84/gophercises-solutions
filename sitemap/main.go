package main

import (
	"flag"
	"fmt"
	"gophercises/sitemap/crawler"
	"gophercises/sitemap/sitemap"
	"log"
)

func main() {
	var url string
	var depth uint

	flag.StringVar(&url, "url", "http://www.ansa.it", "URL of the website to map")
	flag.UintVar(&depth, "depth", 3, "Maximum depth when following links")
	flag.Parse()

	paths, err := crawler.Crawl(url, depth)
	if err != nil {
		log.Fatal(err)
	}

	xmlDoc, err := sitemap.Sitemap(url, paths)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(xmlDoc))
}
