package cmd

import (
	"fmt"
	"log"

	"github.com/SirNoob97/gophercises/task-list/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			log.Fatal("Something went wrong", err.Error())
		}

		if len(tasks) == 0 {
			fmt.Println("You have no task to complete!!")
			return
		}

		fmt.Println("You have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("%d) %s\n", i+1, task.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
