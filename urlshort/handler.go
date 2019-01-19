package urlshort

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if location, ok := pathsToUrls[r.URL.Path]; ok {
			w.Header().Set("Location", location)
			w.WriteHeader(302)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// DataHandler will parse the provided data, in YAML or JSON format,
// following format parameter, and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the data file, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// JSON is expected to be in the format:
//
//  [
//  	{
//  		"path": "/some-path",
//  		"url": "https://www.some-url.com/demo"
//  	}
//  ]
//
// The only errors that can be returned all related to having
// invalid YAML or JSON data.
func DataHandler(data []byte, format string, fallback http.Handler) (http.HandlerFunc, error) {
	locations, err := parseData(data, format)

	return func(w http.ResponseWriter, r *http.Request) {
		requestURL := r.URL.Path

		for url, path := range locations {
			if requestURL == url {
				w.Header().Set("Location", path)
				w.WriteHeader(302)
				return
			}
		}

		fallback.ServeHTTP(w, r)
	}, err
}

func parseData(data []byte, format string) (map[string]string, error) {
	var parsedData []map[string]string
	var err error

	switch format {
	case "yaml":
		err = yaml.Unmarshal(data, &parsedData)
	case "json":
		err = json.Unmarshal(data, &parsedData)
	default:
		err = errors.New("unsupported data format")
	}

	if err != nil {
		return map[string]string{}, err
	}

	locations := make(map[string]string)
	for _, row := range parsedData {
		locations[row["path"]] = row["url"]
	}

	return locations, err
}
