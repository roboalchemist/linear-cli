package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readmeContents string
var skillContents string
var skillRefContents string

// SetReadmeContents sets the README content from main package
func SetReadmeContents(content string) {
	readmeContents = content
}

// SetSkillContents sets the skill content from main package
func SetSkillContents(content, refContent string) {
	skillContents = content
	skillRefContents = refContent
}

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Display the linear-cli documentation",
	Long: `Display the complete linear-cli documentation including README and Claude Code skill.

This command outputs the full documentation in markdown format,
which can be piped to other tools or saved to a file.

Examples:
  linear-cli docs                    # Display documentation
  linear-cli docs | less            # View with pager
  linear-cli docs > linear-cli-docs.md  # Save to file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(readmeContents)
		fmt.Print("\n\n---\n\n# Claude Code Skill\n\n")
		fmt.Print(skillContents)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
