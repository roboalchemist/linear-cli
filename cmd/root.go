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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linctl",
	Short: "A comprehensive Linear CLI tool",
	Long: color.New(color.FgCyan).Sprint(`
┌─────────────────────────────────────────┐
│              🚀 linctl                   │
│      Linear CLI - Built with ❤️         │
└─────────────────────────────────────────┘

A comprehensive CLI tool for Linear's API featuring:
• Issue management (create, list, update, archive)
• Project tracking and collaboration  
• Team and user management
• Comments and attachments
• Webhook configuration
• Table/plaintext/JSON output formats
`),
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.linctl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&plaintext, "plaintext", "p", false, "plaintext output (non-interactive)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "JSON output")

	// Bind flags to viper
	viper.BindPFlag("plaintext", rootCmd.PersistentFlags().Lookup("plaintext"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
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

		// Search config in home directory with name ".linctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".linctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if !plaintext && !jsonOut {
			fmt.Fprintln(os.Stderr, color.New(color.FgGreen).Sprintf("✅ Using config file: %s", viper.ConfigFileUsed()))
		}
	}
}