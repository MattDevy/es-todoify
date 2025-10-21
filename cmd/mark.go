/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/MattDevy/es-todoify/internal/todo"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// markCmd represents the mark command
var markCmd = &cobra.Command{
	Use:     "mark [todo-id]",
	Aliases: []string{"m"},
	Short:   "Mark a todo with a new status",
	Long: `Change the status of an existing todo item.

You can mark a todo with any of the following statuses:
  - pending
  - in_progress
  - completed
  - cancelled
  - blocked

Note: Business rules are enforced. For example, completed todos cannot be marked as blocked.

Examples:
  # Mark a todo as in progress
  todoify mark <uuid> --status in_progress
  todoify m <uuid> -s in_progress

  # Mark a todo as completed
  todoify mark <uuid> --status completed

  # Mark a todo as blocked
  todoify m <uuid> -s blocked`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse and validate UUID
		id, err := uuid.Parse(args[0])
		if err != nil {
			logger.Error("failed to parse todo id", "error", err)
			os.Exit(1)
		}

		// Get and validate status
		statusStr := viper.GetString("status")
		if statusStr == "" {
			logger.Error("status flag is required (use --status or -s)")
			os.Exit(1)
		}

		status := todo.Status(statusStr)
		if !status.IsValid() {
			logger.Error("invalid status", "status", statusStr, "valid_statuses", "pending, in_progress, completed, cancelled, blocked")
			os.Exit(1)
		}

		// Call service to change status
		updatedTodo, err := service.ChangeStatus(cmd.Context(), id.String(), status)
		if err != nil {
			logger.Error("failed to change status", "error", err)
			os.Exit(1)
		}

		// Log success
		logger.Info("successfully changed status", "id", updatedTodo.ID, "status", updatedTodo.Status, "title", updatedTodo.Title)
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Define required status flag
	markCmd.Flags().StringP("status", "s", "", "New status (pending, in_progress, completed, cancelled, blocked)")
	cobra.CheckErr(markCmd.MarkFlagRequired("status"))

	// Bind flag to viper
	viper.BindPFlag("status", markCmd.Flags().Lookup("status"))
}
