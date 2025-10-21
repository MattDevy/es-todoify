/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package operations

import (
	"os"

	"github.com/MattDevy/es-todoify/cmd/sdk"
	esrepo "github.com/MattDevy/es-todoify/internal/todo/repositories/elasticsearch/v9"
	"github.com/spf13/cobra"
)

// NewMigrateCmd creates the migrate command with injected dependencies.
// This command creates or updates Elasticsearch indices based on the defined mappings.
func NewMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Create or update Elasticsearch indices",
		Long: `Create or update Elasticsearch indices based on the defined mappings.

This command will create the necessary indices in Elasticsearch for storing todos.
If the indices already exist, this command will fail. Use with caution in production.

Examples:
  # Create indices
  todoify operations migrate`,
		Run: func(cmd *cobra.Command, args []string) {
			repo := sdk.GetRepo(cmd.Context())
			logger := sdk.GetLogger(cmd.Context())
			// Type assert to Elasticsearch repository
			esRepository, ok := repo.(*esrepo.Repository)
			if !ok {
				logger.Error("migrate command only supports Elasticsearch repository")
				os.Exit(1)
			}

			// Create indices
			err := esRepository.CreateIndices(cmd.Context())
			if err != nil {
				logger.Error("failed to create indices", "error", err)
				os.Exit(1)
			}

			logger.Info("indices created successfully")
		},
	}

	return cmd
}
