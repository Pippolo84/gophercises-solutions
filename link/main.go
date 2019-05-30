package main

import (
	"bufio"
	"flag"
	"fmt"
	"gophercises/link/linkextract"
	"log"
	"os"
)

func main() {
	var fileName string

	flag.StringVar(&fileName, "file", "input.html", "HTML file to parse")
	flag.Parse()

	in, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(in)

	links, err := linkextract.Links(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range links {
		fmt.Println(link)
	}
}
