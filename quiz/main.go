package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type csvLine struct {
	question string
	answer   string
}

func main() {
	var fileName string
	flag.StringVar(
		&fileName,
		"csv",
		"problems.csv",
		"a csv file in the format 'question,answer'",
	)
	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = 2

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	stdinReader := bufio.NewReader(os.Stdin)
	total := len(records)
	score := 0
	for i, r := range records {
		fmt.Printf("Problem #%v: %s = ", i+1, r[0])

		line, err := stdinReader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}
		line = strings.Replace(line, "\r\n", "", -1)
		line = strings.Replace(line, "\n", "", -1)

		if line == r[1] {
			score++
		}
	}

	fmt.Printf("You scored %v out of %v.\n", score, total)
}
