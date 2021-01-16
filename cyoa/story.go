package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

// Story format
type Story map[string]Chapter

// Chapter of the user story
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option of the chapters
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"chapter"`
}

// JSONStory return a Story struct with the json file content
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)

	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

var defaultHandlerTemplate = `
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
            <li> <a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>`

var tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))

// HandlerOption provider of options
type HandlerOption func(h *handler)

// WithTemplate set template to use
func WithTemplate(t *template.Template) HandlerOption {
	// provisional
	//if t == nil {
	//t = tpl
	//}
	return func(h *handler) {
		h.t = t
	}
}

// WithPathFunc set path to go
func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFunc = fn
	}
}

// NewHandler custom handler provider
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type handler struct {
	s        Story
	t        *template.Template
	pathFunc func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Fatalf("Error trying to build the template\n%v", err.Error())
			http.Error(w, "Error trying to build the template", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}
