package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// skillCmd represents the skill command group
var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage the Claude Code skill for linear-cli",
	Long:  `Commands for viewing and installing the embedded Claude Code skill.`,
}

// skillPrintCmd prints the skill content to stdout
var skillPrintCmd = &cobra.Command{
	Use:   "print",
	Short: "Print the embedded Claude Code skill",
	Long: `Print the SKILL.md content to stdout.

Examples:
  linear-cli skill print              # Display skill
  linear-cli skill print | less       # View with pager`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(skillContents)
	},
}

// skillAddCmd installs the skill to ~/.claude/skills/linear-cli/
var skillAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Install the Claude Code skill to ~/.claude/skills/linear-cli/",
	Long: `Copy the embedded skill files to ~/.claude/skills/linear-cli/.

This installs:
  ~/.claude/skills/linear-cli/SKILL.md
  ~/.claude/skills/linear-cli/reference/commands.md

Examples:
  linear-cli skill add    # Install skill files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine home directory: %w", err)
		}

		skillDir := filepath.Join(home, ".claude", "skills", "linear-cli")
		refDir := filepath.Join(skillDir, "reference")

		if err := os.MkdirAll(refDir, 0o755); err != nil {
			return fmt.Errorf("could not create directory %s: %w", refDir, err)
		}

		skillPath := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(skillPath, []byte(skillContents), 0o644); err != nil {
			return fmt.Errorf("could not write %s: %w", skillPath, err)
		}
		fmt.Printf("Wrote %s\n", skillPath)

		refPath := filepath.Join(refDir, "commands.md")
		if err := os.WriteFile(refPath, []byte(skillRefContents), 0o644); err != nil {
			return fmt.Errorf("could not write %s: %w", refPath, err)
		}
		fmt.Printf("Wrote %s\n", refPath)

		fmt.Println("Claude Code skill installed successfully.")
		return nil
	},
}

func init() {
	skillCmd.AddCommand(skillPrintCmd)
	skillCmd.AddCommand(skillAddCmd)
	rootCmd.AddCommand(skillCmd)
}
