/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MattDevy/es-todoify/internal/todo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List todos with optional filtering, sorting, and pagination",
	Long: `List todos from Elasticsearch with powerful filtering and search capabilities.

You can filter by status, labels, search text, and date ranges. Results can be
sorted by different fields and paginated for large result sets.

Examples:
  # List all todos (default: 50 most recent)
  todoify list

  # Filter by status
  todoify list --status pending

  # Filter by multiple labels
  todoify list --labels bug,urgent

  # Full-text search
  todoify list --search "authentication"

  # Pagination
  todoify list --limit 10 --offset 20

  # Custom sorting
  todoify list --sort-by title --sort-order asc

  # Combined filters
  todoify list --status in_progress --labels backend --limit 25`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build filter from flags
		filter, err := buildFilterFromFlags()
		if err != nil {
			logger.Error("invalid filter parameters", "error", err)
			os.Exit(1)
		}

		// Call service to list todos
		todos, err := service.ListTodos(cmd.Context(), filter)
		if err != nil {
			logger.Error("failed to list todos", "error", err)
			os.Exit(1)
		}

		// Print results
		printTodos(todos)
	},
}

// buildFilterFromFlags constructs a ListFilter from command flags
func buildFilterFromFlags() (todo.ListFilter, error) {
	filter := todo.DefaultListFilter()

	// Status filter
	if viper.IsSet("status") {
		statusStr := viper.GetString("status")
		status := todo.Status(statusStr)
		if !status.IsValid() {
			return filter, fmt.Errorf("invalid status: %s (valid: pending, in_progress, completed, cancelled, blocked)", statusStr)
		}
		filter.Status = status
	}

	// Labels filter
	if viper.IsSet("labels") {
		filter.Labels = viper.GetStringSlice("labels")
	}

	// Search query
	if viper.IsSet("search") {
		filter.SearchQuery = viper.GetString("search")
	}

	// From date filter
	if viper.IsSet("from-date") {
		fromDateStr := viper.GetString("from-date")
		fromDate, err := time.Parse(time.RFC3339, fromDateStr)
		if err != nil {
			return filter, fmt.Errorf("invalid from-date format (use RFC3339, e.g., 2025-01-15T00:00:00Z): %w", err)
		}
		filter.FromDate = &fromDate
	}

	// To date filter
	if viper.IsSet("to-date") {
		toDateStr := viper.GetString("to-date")
		toDate, err := time.Parse(time.RFC3339, toDateStr)
		if err != nil {
			return filter, fmt.Errorf("invalid to-date format (use RFC3339, e.g., 2025-01-15T23:59:59Z): %w", err)
		}
		filter.ToDate = &toDate
	}

	// Pagination
	if viper.IsSet("limit") {
		filter.Limit = viper.GetInt("limit")
	}
	if viper.IsSet("offset") {
		filter.Offset = viper.GetInt("offset")
	}

	// Sorting
	if viper.IsSet("sort-by") {
		sortByStr := viper.GetString("sort-by")
		sortBy := todo.SortField(sortByStr)
		if !sortBy.IsValid() {
			return filter, fmt.Errorf("invalid sort-by: %s (valid: createTime, updateTime, title, status)", sortByStr)
		}
		filter.SortBy = sortBy
	}

	if viper.IsSet("sort-order") {
		sortOrderStr := viper.GetString("sort-order")
		sortOrder := todo.SortOrder(sortOrderStr)
		if !sortOrder.IsValid() {
			return filter, fmt.Errorf("invalid sort-order: %s (valid: asc, desc)", sortOrderStr)
		}
		filter.SortOrder = sortOrder
	}

	return filter, nil
}

// printTodos prints todos in a simple, readable format
func printTodos(todos []*todo.Todo) {
	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}

	fmt.Printf("Found %d todo(s):\n\n", len(todos))

	for i, t := range todos {
		fmt.Printf("ID:          %s\n", t.ID)
		fmt.Printf("Title:       %s\n", t.Title)
		if t.Description != "" {
			fmt.Printf("Description: %s\n", t.Description)
		}
		fmt.Printf("Status:      %s\n", t.Status)
		if len(t.Labels) > 0 {
			fmt.Printf("Labels:      [%s]\n", strings.Join(t.Labels, ", "))
		}
		fmt.Printf("Created:     %s\n", t.CreateTime.Format(time.RFC3339))
		fmt.Printf("Updated:     %s\n", t.UpdateTime.Format(time.RFC3339))

		// Add separator between todos (except after the last one)
		if i < len(todos)-1 {
			fmt.Println("\n---")
			fmt.Println()
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Filter flags
	listCmd.Flags().StringP("status", "s", "", "Filter by status (pending, in_progress, completed, cancelled, blocked)")
	listCmd.Flags().StringSliceP("labels", "l", []string{}, "Filter by labels (comma-separated, must have all)")
	listCmd.Flags().StringP("search", "q", "", "Search query for title and description")
	listCmd.Flags().String("from-date", "", "Filter todos created on or after this date (RFC3339 format)")
	listCmd.Flags().String("to-date", "", "Filter todos created on or before this date (RFC3339 format)")

	// Pagination flags
	listCmd.Flags().Int("limit", 50, "Maximum number of results to return")
	listCmd.Flags().Int("offset", 0, "Number of results to skip (for pagination)")

	// Sorting flags
	listCmd.Flags().String("sort-by", "createTime", "Field to sort by (createTime, updateTime, title, status)")
	listCmd.Flags().String("sort-order", "desc", "Sort order (asc, desc)")

	// Bind flags to viper
	viper.BindPFlag("status", listCmd.Flags().Lookup("status"))
	viper.BindPFlag("labels", listCmd.Flags().Lookup("labels"))
	viper.BindPFlag("search", listCmd.Flags().Lookup("search"))
	viper.BindPFlag("from-date", listCmd.Flags().Lookup("from-date"))
	viper.BindPFlag("to-date", listCmd.Flags().Lookup("to-date"))
	viper.BindPFlag("limit", listCmd.Flags().Lookup("limit"))
	viper.BindPFlag("offset", listCmd.Flags().Lookup("offset"))
	viper.BindPFlag("sort-by", listCmd.Flags().Lookup("sort-by"))
	viper.BindPFlag("sort-order", listCmd.Flags().Lookup("sort-order"))
}
