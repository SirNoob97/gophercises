package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	fileName := "birthday_001.txt"
	newName, err := match(fileName, 4)
	if err != nil {
		log.Fatal("Not Match")
	}
	fmt.Println(newName)
}

func match(fileName string, total int) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match the pattern", fileName)
	}

	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
