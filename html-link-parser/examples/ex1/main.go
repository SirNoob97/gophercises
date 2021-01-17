package main

import (
	"fmt"
	"os"

	link "github.com/SirNoob97/gophercises/html-link-parser"
)

var ex1HTML = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">
    A link to another page 
    <span> this is a span</span>
  </a>
  <a href="/page-two">A link to another page</a>
</body>
</html>`

func main() {
	// create a <Reader> from a string
	//r := strings.NewReader(ex1HTML)

	r, err := os.Open("../../ex1.html")
	links, err := link.Parser(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", links)
}
