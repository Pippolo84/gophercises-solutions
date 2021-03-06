package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type csvLine struct {
	question string
	answer   string
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	var fileName string
	flag.StringVar(
		&fileName,
		"csv",
		"problems.csv",
		"a csv file in the format 'question,answer'",
	)

	var limit int
	flag.IntVar(
		&limit,
		"limit",
		30,
		"the time limit for the quiz in seconds",
	)

	var randomize bool
	flag.BoolVar(
		&randomize,
		"randomize",
		false,
		"randomize the order of the questions",
	)

	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = 2
	csvReader.LazyQuotes = true

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	if randomize {
		shuffle(records)
	}

	score := quiz(bufio.NewReader(os.Stdin), records, limit)

	fmt.Printf("You scored %v out of %v.\n", score, len(records))
}

func quiz(reader *bufio.Reader, records [][]string, limit int) (correct int) {
	answers := make(chan string)
	defer close(answers)

	timeout := time.After(time.Duration(limit) * time.Second)
	for i, r := range records {
		fmt.Printf("Problem #%v: %s = ", i+1, r[0])

		go readInput(reader, answers)

		select {
		case answer := <-answers:
			if cleanString(answer) == cleanString(r[1]) {
				correct++
			}
		case <-timeout:
			fmt.Println()
			return
		}
	}

	return
}

func readInput(r *bufio.Reader, in chan<- string) {
	answer, err := r.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}

	answer = strings.Replace(answer, "\r\n", "", -1)
	answer = strings.Replace(answer, "\n", "", -1)

	in <- answer
}

func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func shuffle(s [][]string) {
	for n := len(s); n > 0; n-- {
		randIndex := rand.Intn(n)
		for i := 0; i < len(s[0]); i++ {
			s[n-1][i], s[randIndex][i] = s[randIndex][i], s[n-1][i]
		}
	}
}
