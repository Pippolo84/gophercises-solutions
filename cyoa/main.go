package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

// AppStatus stores the global application status
type AppStatus struct {
	Arcs map[string]Arc
	Tmpl *template.Template
}

var appStatus AppStatus

// Arc stores info about a story arc
type Arc struct {
	Title     string   `json:"title"`
	Story     []string `json:"story"`
	Options   []Option `json:"options"`
	HasOption bool
}

// Option contains info to jump from one arc to another
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {
	jsonFile, err := os.Open("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	arcs := make(map[string]Arc)
	err = json.Unmarshal(byteValue, &arcs)
	if err != nil {
		log.Fatal(err)
	}

	key := "intro"
	for {
		arc, ok := arcs[key]
		if !ok {
			log.Fatalf("unknown key %s", key)
		}

		fmt.Printf("\n%s\n\n", arc.Title)

		for _, paragraph := range arc.Story {
			fmt.Println(paragraph)
		}

		if len(arc.Options) == 0 {
			fmt.Printf("\nThe End\n")
			break
		}

		for {
			fmt.Println()

			for idx, option := range arc.Options {
				fmt.Printf("%d) %s\n", idx+1, option.Text)
			}
			fmt.Printf("\nYour choice: ")

			var choice int
			_, err := fmt.Scanf("%d", &choice)
			if err == nil && choice > 0 && choice <= len(arc.Options) {
				key = arcs[key].Options[choice-1].Arc
				break
			}

			fmt.Printf("\ninvalid choice, please try again\n")
		}
	}
}
