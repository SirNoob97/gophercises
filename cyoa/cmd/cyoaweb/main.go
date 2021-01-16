package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SirNoob97/gophercises/cyoa"
)

func main() {
	port := flag.Int("port", 8080, "the port to start the CYOA web application")
	file := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the JSON file: %s", *file))
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the JSON file: %s", *file))
	}

	fmt.Printf("Using the story in %s\n", *file)

	tpl := template.Must(template.New("").Parse(defaultTemplate))

	//h := cyoa.NewHandler(story, cyoa.WithTemplate(nil))
	h := cyoa.NewHandler(story,
		cyoa.WithTemplate(tpl),
		cyoa.WithPathFunc(pathFn),
	)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("Starting the server on :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	return path[len("/story/"):]
}

func exit(msg string) {
	log.Fatalln(msg)
}

var defaultTemplate = `
<!DOCTYPE html>
<head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
    <ul>
        {{range .Options}}
            <li> <a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>`
