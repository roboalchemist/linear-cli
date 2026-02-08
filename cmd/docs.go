package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readmeContents string

// SetReadmeContents sets the README content from main package
func SetReadmeContents(content string) {
	readmeContents = content
}

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Display the linear-cli documentation",
	Long: `Display the complete linear-cli documentation from README.md.

This command outputs the full documentation in markdown format,
which can be piped to other tools or saved to a file.

Examples:
  linear-cli docs                    # Display documentation
  linear-cli docs | less            # View with pager
  linear-cli docs > linear-cli-docs.md  # Save to file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(readmeContents)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
