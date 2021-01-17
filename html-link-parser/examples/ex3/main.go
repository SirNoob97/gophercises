package main

import (
	"fmt"
	"os"

	link "github.com/SirNoob97/gophercises/html-link-parser"
)

func main() {
	// create a <Reader> from a string
	//r := strings.NewReader(ex1HTML)

	r, err := os.Open("../../ex3.html")
	links, err := link.Parser(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", links)
}
