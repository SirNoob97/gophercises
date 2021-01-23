package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type file struct {
	name string
	path string
}

type matchResult struct {
	base  string
	index int
	ext   string
}

var pattern = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	walkDir := "sample"
	var toRename []string
	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}
		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		filename := filepath.Base(oldPath)
		newFileName, _ := match(filename)
		newPath := filepath.Join(dir, newFileName)
		fmt.Printf("mv %s => %s\n", oldPath, newPath)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("Error renaming: %s to %s. %s\n", oldPath, newPath, err.Error())
		}
	}
}

func match(fileName string) (string, error) {
	if !pattern.MatchString(fileName) {
		return "", fmt.Errorf("%s didn't match the pattern", fileName)
	}
	return pattern.ReplaceAllString(fileName, replaceString), nil
}
