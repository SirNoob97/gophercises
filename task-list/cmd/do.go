package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/SirNoob97/gophercises/task-list/db"
	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete.",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument:", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := db.AllTasks()
		if err != nil {
			log.Fatalln("Something went wrong", err)
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number:", id)
				continue
			}
			task := tasks[id-1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				log.Fatalf("Failed to mark \"%d\" task as completed. Error %s", id, err)
			} else {
				fmt.Printf("Mark \"%d\" task as completed.\n", id)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
