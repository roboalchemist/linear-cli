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
	Short: "Display the linctl documentation",
	Long: `Display the complete linctl documentation from README.md.

This command outputs the full documentation in markdown format,
which can be piped to other tools or saved to a file.

Examples:
  linctl docs                    # Display documentation
  linctl docs | less            # View with pager
  linctl docs > linctl-docs.md  # Save to file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(readmeContents)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
