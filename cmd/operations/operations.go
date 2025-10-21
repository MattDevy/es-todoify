/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package operations

import (
	"github.com/spf13/cobra"
)

// Register registers the operations command and all subcommands to the parent command.
// This function creates the complete operations command hierarchy with injected dependencies.
func Register(parentCmd *cobra.Command) {
	// Create operations parent command
	opsCmd := NewOperationsCmd()

	// Add subcommands with dependencies
	opsCmd.AddCommand(NewMigrateCmd())
	opsCmd.AddCommand(NewHealthCmd())

	// Register operations command to parent
	parentCmd.AddCommand(opsCmd)
}

// NewOperationsCmd creates the operations parent command.
// This command groups administrative and maintenance operations.
func NewOperationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operations",
		Aliases: []string{"ops"},
		Short:   "Administrative and maintenance operations",
		Long: `Administrative and maintenance operations for managing the Elasticsearch backend.

This includes operations like creating indices, running migrations, and other
administrative tasks that are separate from day-to-day todo management.`,
	}

	return cmd
}
