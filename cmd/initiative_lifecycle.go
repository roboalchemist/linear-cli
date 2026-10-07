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

var initiativeArchiveCmd = &cobra.Command{
	Use:     "archive INITIATIVE-ID",
	Aliases: []string{"rm"},
	Short:   "Archive an initiative",
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
		if err := client.ArchiveInitiative(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to archive initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Archived initiative %s", args[0]), plaintext, jsonOut)
	},
}

var initiativeUnarchiveCmd = &cobra.Command{
	Use:   "unarchive INITIATIVE-ID",
	Short: "Restore an archived initiative",
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
		if err := client.UnarchiveInitiative(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to unarchive initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Unarchived initiative %s", args[0]), plaintext, jsonOut)
	},
}

var initiativeLabelAddCmd = &cobra.Command{
	Use:   "add INITIATIVE-ID LABEL-ID [LABEL-ID...]",
	Short: "Add label(s) to an initiative",
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
			if err := client.AddInitiativeLabel(context.Background(), args[0], labelID); err != nil {
				output.Error(fmt.Sprintf("Failed to add label %s: %v", labelID, err), plaintext, jsonOut)
				os.Exit(1)
			}
		}
		output.Success(fmt.Sprintf("Added %d label(s) to %s", len(args)-1, args[0]), plaintext, jsonOut)
	},
}

var initiativeLabelRemoveCmd = &cobra.Command{
	Use:     "remove INITIATIVE-ID LABEL-ID [LABEL-ID...]",
	Aliases: []string{"rm"},
	Short:   "Remove label(s) from an initiative",
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
			if err := client.RemoveInitiativeLabel(context.Background(), args[0], labelID); err != nil {
				output.Error(fmt.Sprintf("Failed to remove label %s: %v", labelID, err), plaintext, jsonOut)
				os.Exit(1)
			}
		}
		output.Success(fmt.Sprintf("Removed %d label(s) from %s", len(args)-1, args[0]), plaintext, jsonOut)
	},
}

func init() {
	initiativeCmd.AddCommand(initiativeArchiveCmd)
	initiativeCmd.AddCommand(initiativeUnarchiveCmd)

	// Extend the existing `initiative label` command (defined in initiative_label.go).
	initiativeLabelCmd.AddCommand(initiativeLabelAddCmd)
	initiativeLabelCmd.AddCommand(initiativeLabelRemoveCmd)
}
