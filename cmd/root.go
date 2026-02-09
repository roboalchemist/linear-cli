package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	plaintext bool
	jsonOut   bool
)

// version is set at build time via -ldflags
// default value is for local dev builds
var version = "dev"

// generateHeader creates a nice header box with proper Unicode box drawing
func generateHeader() string {
	return "" +
		"┌───────────────────────────┐\n" +
		"│       linear-cli          │\n" +
		"│   CLI for the Linear API  │\n" +
		"└───────────────────────────┘"
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "linear-cli",
	Short:   "A comprehensive Linear CLI tool",
	Long:    color.New(color.FgCyan).Sprintf("%s\nA CLI for Linear's API featuring:\n• Issues, projects, cycles, labels, documents, initiatives, views\n• Comments, attachments, relations, milestones, status updates\n• Team and user management\n• Raw GraphQL queries\n• Table/plaintext/JSON output formats\n", generateHeader()),
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetRootCmd returns the root command for testing
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.linear-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&plaintext, "plaintext", "p", false, "plaintext output (non-interactive)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "JSON output")

	// Bind flags to viper
	_ = viper.BindPFlag("plaintext", rootCmd.PersistentFlags().Lookup("plaintext"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".linear-cli" (without extension).
		// Also check legacy ".linctl" config for backward compat.
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".linear-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if !plaintext && !jsonOut {
			fmt.Fprintln(os.Stderr, color.New(color.FgGreen).Sprintf("✅ Using config file: %s", viper.ConfigFileUsed()))
		}
	}
}
