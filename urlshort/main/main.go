package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/mediocregopher/radix.v2/redis"
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

	pathsToUrls, err := getLocationsMap()
	if err != nil {
		panic(err)
	}

	// Build the defaultHandler using the mux as the fallback
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
	fmt.Fprintf(w, "No mapping defined for this handler :(")
}

func getLocationsMap() (map[string]string, error) {
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp := client.Cmd("LLEN", "urlshort")
	if resp.Err != nil {
		return nil, resp.Err
	}

	l, err := resp.Int()
	if err != nil {
		return nil, err
	}

	resp = client.Cmd("LRANGE", "urlshort", "0", strconv.Itoa(l))
	if resp.Err != nil {
		return nil, resp.Err
	}

	pathsToUrls, err := resp.Map()
	if err != nil {
		return nil, err
	}

	return pathsToUrls, nil
}
