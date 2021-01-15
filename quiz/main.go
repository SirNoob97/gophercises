package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "a csv files in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFile))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the CSV file")
	}

	problems := parseLines(lines)

	time := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	count := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCn := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCn <- answer
		}()
		select {
		case <-time.C:
			fmt.Printf("\nYou scored %d out of %d.\n", count, len(problems))
			return
		case answer := <-answerCn:
			if answer == p.a {
				count++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", count, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	log.Fatalln(msg)
}
