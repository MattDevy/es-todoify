/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package operations

import (
	"encoding/json"
	"os"

	"github.com/MattDevy/es-todoify/cmd/sdk"
	"github.com/spf13/cobra"
)

// NewHealthCmd creates the health command with injected dependencies.
// This command checks the health of the Elasticsearch backend and displays comprehensive health information.
func NewHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check the health of the Elasticsearch backend",
		Long: `Check the health of the Elasticsearch backend and display comprehensive health information.

This command performs a health check on the Elasticsearch cluster and returns information including:
- Overall health status (healthy/degraded/unhealthy)
- Cluster availability
- Response time
- Node count
- Elasticsearch version
- Detailed cluster statistics (shards, pending tasks, etc.)

Examples:
  # Check backend health
  todoify operations health
  
  # Check backend health (short alias)
  todoify ops health`,
		Run: func(cmd *cobra.Command, args []string) {
			service := sdk.GetService(cmd.Context())
			logger := sdk.GetLogger(cmd.Context())

			// Perform health check
			healthInfo, err := service.Health(cmd.Context())
			if err != nil {
				logger.Error("health check failed", "error", err)

				// If we have partial health info, still display it
				if healthInfo != nil {
					outputJSON(healthInfo, logger)
				}
				os.Exit(1)
			}

			// Display health information as formatted JSON
			outputJSON(healthInfo, logger)

			// Exit with non-zero code if unhealthy
			if healthInfo.Status != "healthy" {
				os.Exit(1)
			}
		},
	}

	return cmd
}

// outputJSON marshals and prints the health info as formatted JSON
func outputJSON(healthInfo interface{}, logger interface{}) {
	jsonBytes, err := json.MarshalIndent(healthInfo, "", "  ")
	if err != nil {
		if l, ok := logger.(interface{ Error(string, ...interface{}) }); ok {
			l.Error("failed to marshal health info", "error", err)
		}
		return
	}

	// Print to stdout
	os.Stdout.Write(jsonBytes)
	os.Stdout.Write([]byte("\n"))
}
