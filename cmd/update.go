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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update [todo-id]",
	Aliases: []string{"u"},
	Short:   "Update a todo's title, description, or labels",
	Long: `Update one or more fields of an existing todo item.

You can update the title, description, and/or labels of a todo by providing
its UUID and one or more update flags. At least one field must be provided.

Examples:
  # Update just the title
  todoify update abc123-... --title "New title"

  # Update multiple fields
  todoify update abc123-... -t "New title" -d "Updated description"

  # Update only labels
  todoify update abc123-... --labels bug,urgent,backend`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse and validate UUID
		id, err := uuid.Parse(args[0])
		if err != nil {
			logger.Error("failed to parse todo id", "error", err)
			os.Exit(1)
		}

		// Build UpdateTodo struct from provided flags
		update := buildUpdateFromFlags()

		// Check if at least one field is provided
		if update.Title == nil && update.Description == nil && update.Labels == nil {
			logger.Error("at least one field must be provided to update (--title, --description, or --labels)")
			os.Exit(1)
		}

		// Call service to update todo
		updatedTodo, err := service.UpdateTodo(cmd.Context(), id.String(), update)
		if err != nil {
			// Check for validation errors and translate them
			if validationErr := todo.TranslateError(err); validationErr != nil {
				logger.Error("validation failed", "error", validationErr)
			} else {
				logger.Error("failed to update todo", "error", err)
			}
			os.Exit(1)
		}

		// Log success
		logger.Info("successfully updated todo", "id", updatedTodo.ID, "title", updatedTodo.Title)
	},
}

// buildUpdateFromFlags constructs an UpdateTodo struct from command flags.
// Only fields that were explicitly provided via flags are set (using pointers).
func buildUpdateFromFlags() todo.UpdateTodo {
	update := todo.UpdateTodo{}

	// Check if title flag was provided
	if viper.IsSet("title") {
		title := viper.GetString("title")
		update.Title = &title
	}

	// Check if description flag was provided
	if viper.IsSet("description") {
		description := viper.GetString("description")
		update.Description = &description
	}

	// Check if labels flag was provided
	if viper.IsSet("labels") {
		labels := viper.GetStringSlice("labels")
		update.Labels = labels
	}

	return update
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Define flags for editable fields
	updateCmd.Flags().StringP("title", "t", "", "New title for the todo")
	updateCmd.Flags().StringP("description", "d", "", "New description for the todo")
	updateCmd.Flags().StringSliceP("labels", "l", []string{}, "New labels for the todo (comma-separated)")

	// Bind flags to viper so we can check if they were set
	viper.BindPFlag("title", updateCmd.Flags().Lookup("title"))
	viper.BindPFlag("description", updateCmd.Flags().Lookup("description"))
	viper.BindPFlag("labels", updateCmd.Flags().Lookup("labels"))
}
