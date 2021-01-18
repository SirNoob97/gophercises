package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("do called")
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
