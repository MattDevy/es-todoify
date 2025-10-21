/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete [todo-id]",
	Aliases: []string{"d"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := uuid.Parse(args[0])
		if err != nil {
			logger.Error("failed to parse todo id", "error", err)
			os.Exit(1)
		}
		err = service.DeleteTodo(cmd.Context(), id.String())
		if err != nil {
			logger.Error("failed to delete todo", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
