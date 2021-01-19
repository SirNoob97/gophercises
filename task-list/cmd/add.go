package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/SirNoob97/gophercises/task-list/db"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the list.",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")

		i, err := db.CreateTask(task)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if i > 0 {
			fmt.Printf("Add \"%s\" to your task list.\n", task)
		} else {
			log.Fatal("Something went wrong.")
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
