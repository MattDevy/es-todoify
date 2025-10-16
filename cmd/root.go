/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	esClient *elasticsearch.TypedClient // Typed Elasticsearch client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todoify",
	Short: "A CLI todo application powered by Elasticsearch",
	Long: `Todoify is a command-line todo application that uses Elasticsearch as its backend.

It provides a simple and powerful way to manage your todos with features like:
- Create, update, list, and delete todos
- Search and filter todos
- Bulk operations with JSON and CSV support
- Todo statistics and insights

All data is stored in Elasticsearch, giving you the power of full-text search,
aggregations, and scalability for your todo management.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Viper configuration
		if err := initConfig(cmd); err != nil {
			return err
		}

		// Initialize Elasticsearch client
		if err := initElasticsearchClient(); err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todoify.yaml)")

	// Elasticsearch connection flags
	rootCmd.PersistentFlags().StringSlice("es-addrs", []string{"http://localhost:9200"}, "Elasticsearch addresses (comma-separated)")
	rootCmd.PersistentFlags().String("es-username", "", "Elasticsearch username")
	rootCmd.PersistentFlags().String("es-password", "", "Elasticsearch password")
	rootCmd.PersistentFlags().String("es-api-key", "", "Elasticsearch API key")
	rootCmd.PersistentFlags().String("es-index", "todos", "Elasticsearch index name")

}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".todoify" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".todoify")
	}

	viper.SetEnvPrefix("TODOIFY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// Bind local flags
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// Set defaults (in case not provided anywhere)
	viper.SetDefault("es-addrs", []string{"http://localhost:9200"})
	viper.SetDefault("es-index", "todos")

	return nil
}

// initElasticsearchClient initializes the Elasticsearch typed client
func initElasticsearchClient() error {
	// Retrieve configuration values from Viper
	esAddrs := viper.GetStringSlice("es-addrs")
	esUsername := viper.GetString("es-username")
	esPassword := viper.GetString("es-password")
	esAPIKey := viper.GetString("es-api-key")

	// Validate required configuration
	if len(esAddrs) == 0 {
		return errors.New("at least one Elasticsearch address is required (use --es-addrs, TODOIFY_ES_ADDRS, or config file)")
	}

	// Validate that both username and password are provided together
	if (esUsername != "" && esPassword == "") || (esUsername == "" && esPassword != "") {
		return errors.New("both es-username and es-password must be provided together")
	}

	// Validate that API key and username/password are mutually exclusive
	if esAPIKey != "" && (esUsername != "" || esPassword != "") {
		return errors.New("cannot use both API key and username/password authentication; choose one method")
	}

	// Build Elasticsearch client config
	cfg := elasticsearch.Config{
		Addresses: esAddrs,
	}

	// Add authentication if provided
	if esAPIKey != "" {
		cfg.APIKey = esAPIKey
	} else if esUsername != "" && esPassword != "" {
		cfg.Username = esUsername
		cfg.Password = esPassword
	}

	// Create typed client
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Verify connection by pinging the cluster
	ctx := context.Background()
	info, err := client.Info().Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Elasticsearch at %v: %w\nPlease check:\n  - Elasticsearch is running\n  - Addresses are correct\n  - Credentials are valid", esAddrs, err)
	}

	// Store client in package variable
	esClient = client

	// Log successful connection
	fmt.Printf("Connected to Elasticsearch cluster: %s\n", info.ClusterName)

	return nil
}
