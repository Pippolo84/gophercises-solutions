package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AppStatus struct {
	Arcs map[string]Arc
	Tmpl *template.Template
}

var appStatus AppStatus

type Arc struct {
	Title     string   `json:"title"`
	Story     []string `json:"story"`
	Options   []Option `json:"options"`
	HasOption bool
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type AppHandler func(http.ResponseWriter, *http.Request)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var arc Arc

	if r.URL.Path == "/" {
		arc = appStatus.Arcs["intro"]
	} else {
		var ok bool
		arc, ok = appStatus.Arcs[r.URL.Path[1:]]
		if !ok {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}
	}

	arc.HasOption = len(arc.Options) > 0
	appStatus.Tmpl.Execute(w, arc)
}

func main() {
	jsonFile, err := os.Open("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &appStatus.Arcs)
	if err != nil {
		log.Fatal(err)
	}

	appStatus.Tmpl = template.Must(template.ParseFiles("tmpl/template.html"))

	var appHandler AppHandler

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", appHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
