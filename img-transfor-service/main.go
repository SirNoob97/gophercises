package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SirNoob97/gophercises/img-transfor-service/primitive"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
      <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="image">
        <button type="submit">Upload Image</button>
      </form>
    </body></html>`
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/modify/", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("./img/" + filepath.Base(r.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		ext := filepath.Ext(f.Name())[1:]
		modeStr := r.FormValue("mode")
		if modeStr == "" {
			renderMoreChoices(w, r, f, ext)
			return
		}
		mode, err := strconv.Atoi(modeStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		renderMoreChoices(w, r, f, ext, mode)
		w.Header().Set("Content-Type", fmt.Sprintf("image/%s", ext))
		io.Copy(w, f)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]
		onDisk, err := tempfile("", ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer onDisk.Close()
		_, err = io.Copy(onDisk, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		imgName := filepath.Base(onDisk.Name())
		redirURL := fmt.Sprintf("/modify/%s", imgName)
		http.Redirect(w, r, redirURL, http.StatusFound)
	})
	fs := http.FileServer(http.Dir("./img/"))
	// StripPrefix helps avoid conflicts between url path and directory path
	mux.Handle("/img/", http.StripPrefix("/img", fs))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("Failed to create temporary file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}

func genImage(f io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(f, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}

	outFile, err := tempfile("", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}

func renderMoreChoices(w http.ResponseWriter, r *http.Request, f io.ReadSeeker, ext string, mode ...int) {
	a, err := genImage(f, ext, 10, primitive.ModeCircle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
  f.Seek(0, 0)
	b, err := genImage(f, ext, 10, primitive.ModePolygon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
  f.Seek(0, 0)
	c, err := genImage(f, ext, 10, primitive.ModeRect)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
  f.Seek(0, 0)
	d, err := genImage(f, ext, 10, primitive.ModeTriangle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
  f.Seek(0, 0)

	html := `<html>
  <body>
    {{range .}}
    <a href="/modify/{{.Name}}?mode={{.Mode}}">
      <img style="width: 20%;" src="/img/{{.Name}}">
    </a>
    {{end}}
  </body>
  </html>`

	tmpl := template.Must(template.New("").Parse(html))
	data := []struct {
		Name string
		Mode primitive.Mode
	}{
		{Name: filepath.Base(a), Mode: primitive.ModeCircle},
		{Name: filepath.Base(b), Mode: primitive.ModePolygon},
		{Name: filepath.Base(c), Mode: primitive.ModeRect},
		{Name: filepath.Base(d), Mode: primitive.ModeTriangle},
	}
	tmpl.Execute(w, data)
}
