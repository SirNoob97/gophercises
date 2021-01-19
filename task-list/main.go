package main

import (
	"log"
	"path/filepath"

	"github.com/SirNoob97/gophercises/task-list/cmd"
	"github.com/SirNoob97/gophercises/task-list/db"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		exit(err.Error())
	}

	dbPath := filepath.Join(home, "tasks.db")
	err = db.Init(dbPath)
	if err != nil {
		exit(err.Error())
	}

	cmd.Execute()
}

func exit(msg string) {
	log.Fatal(msg)
}
