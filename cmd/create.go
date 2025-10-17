/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		todo, err := service.CreateTodo(cmd.Context(), viper.GetString("title"), viper.GetString("description"), viper.GetStringSlice("labels"))
		cobra.CheckErr(err)
		logger.Info("created todo", "todo", todo)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("title", "t", "", "The title of the todo")
	cobra.CheckErr(createCmd.MarkFlagRequired("title"))
	createCmd.Flags().StringP("description", "d", "", "The description of the todo")
	createCmd.Flags().StringSliceP("labels", "l", []string{}, "The labels of the todo")
}
