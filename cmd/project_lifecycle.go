package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectUnarchiveCmd = &cobra.Command{
	Use:   "unarchive PROJECT-ID",
	Short: "Restore an archived project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UnarchiveProject(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to unarchive project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Unarchived project %s", args[0]), plaintext, jsonOut)
	},
}

var projectLabelAddCmd = &cobra.Command{
	Use:   "add PROJECT-ID LABEL-ID [LABEL-ID...]",
	Short: "Add label(s) to a project",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		for _, labelID := range args[1:] {
			if err := client.AddProjectLabel(context.Background(), args[0], labelID); err != nil {
				output.Error(fmt.Sprintf("Failed to add label %s: %v", labelID, err), plaintext, jsonOut)
				os.Exit(1)
			}
		}
		output.Success(fmt.Sprintf("Added %d label(s) to %s", len(args)-1, args[0]), plaintext, jsonOut)
	},
}

var projectLabelRemoveCmd = &cobra.Command{
	Use:     "remove PROJECT-ID LABEL-ID [LABEL-ID...]",
	Aliases: []string{"rm"},
	Short:   "Remove label(s) from a project",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		for _, labelID := range args[1:] {
			if err := client.RemoveProjectLabel(context.Background(), args[0], labelID); err != nil {
				output.Error(fmt.Sprintf("Failed to remove label %s: %v", labelID, err), plaintext, jsonOut)
				os.Exit(1)
			}
		}
		output.Success(fmt.Sprintf("Removed %d label(s) from %s", len(args)-1, args[0]), plaintext, jsonOut)
	},
}

var projectReassignStatusCmd = &cobra.Command{
	Use:   "reassign-status",
	Short: "Move projects from one status to another",
	Long:  `Reassign all projects from an original project status to a new one.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		newStatus, _ := cmd.Flags().GetString("new")
		original, _ := cmd.Flags().GetString("original")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.ReassignProjectStatus(context.Background(), newStatus, original); err != nil {
			output.Error(fmt.Sprintf("Failed to reassign project status: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Reassigned projects from %s to %s", original, newStatus), plaintext, jsonOut)
	},
}

var projectStatusUnarchiveCmd = &cobra.Command{
	Use:   "unarchive UPDATE-ID",
	Short: "Restore an archived project status update",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UnarchiveProjectUpdate(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to unarchive project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Unarchived project update %s", args[0]), plaintext, jsonOut)
	},
}

var projectStatusDeleteCmd = &cobra.Command{
	Use:     "purge UPDATE-ID",
	Aliases: []string{"rm-permanent"},
	Short:   "Permanently delete a project status update",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.DeleteProjectUpdate(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Deleted project update %s", args[0]), plaintext, jsonOut)
	},
}

func init() {
	projectCmd.AddCommand(projectUnarchiveCmd)

	// Extend the existing `project label` command (defined in project_label.go).
	projectLabelCmd.AddCommand(projectLabelAddCmd)
	projectLabelCmd.AddCommand(projectLabelRemoveCmd)

	projectCmd.AddCommand(projectReassignStatusCmd)
	projectReassignStatusCmd.Flags().String("new", "", "New project status ID (required)")
	projectReassignStatusCmd.Flags().String("original", "", "Original project status ID (required)")
	_ = projectReassignStatusCmd.MarkFlagRequired("new")
	_ = projectReassignStatusCmd.MarkFlagRequired("original")

	// Extend the existing `project status` command (defined in project_update.go).
	projectStatusCmd.AddCommand(projectStatusUnarchiveCmd)
	projectStatusCmd.AddCommand(projectStatusDeleteCmd)
}
