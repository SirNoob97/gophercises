package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode ...
type Mode int

// Modes ...
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedrect
	ModeBeziers
	ModeRotatedellipse
	ModePolygon
)

// WithMode is an option for the Transform function that will define the mode
// you want to use. By default, ModeTriangle will be used.
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

// Transform will take the provided image and apply a primitive transformation
// to it, then return a reader to the resulting image.
func Transform(image io.Reader, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	in, err := tempfile("in_", "png")
	if err != nil {
		return nil, errors.New("Failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	out, err := tempfile("in_", "png")
	if err != nil {
		return nil, errors.New("Failed to create temporary output file")
	}
	defer os.Remove(out.Name())

	// Read image into "in" file
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("Failed to copy image into temporary input file")
	}

	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to run primitive commad. stdcombo = %s", stdCombo)
	}
	fmt.Println(stdCombo)

	// Read "out" into a reader, return reader, delete "in" and "out" due to defer statement
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("Failed to copy output file into byte buffer")
	}
	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int, args ...string) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	args = append(strings.Fields(argStr), args...)

	cmd := exec.Command("primitive", args...)

	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, errors.New("Failed to create temporary file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
