package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	locationsFile := flag.String("locations", "redirects.yml", "input file name containing the redirects map")
	format := flag.String("format", "yaml", "input file format (yaml|json)")
	flag.Parse()

	locations, err := ioutil.ReadFile(*locationsFile)
	if err != nil {
		log.Fatal(err)
	}

	if *format != "yaml" && *format != "json" {
		log.Fatalf("unsupported input file format: %s\n", *format)
	}

	mux := defaultMux()

	// Build the defaultHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	defaultHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the defaultHandler as the fallback
	handler, err := urlshort.DataHandler(locations, *format, defaultHandler)
	if err != nil {
		panic(err)
	}

	log.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
