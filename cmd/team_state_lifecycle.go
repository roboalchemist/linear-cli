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

var teamStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Manage a team's workflow states",
	Long: `Create, update, or archive a team's workflow states.

State types: backlog, unstarted, started, completed, canceled.`,
}

var teamStateCreateCmd = &cobra.Command{
	Use:   "create --team-id TEAM-ID --name NAME --color HEX --type TYPE",
	Short: "Create a workflow state",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		teamID, _ := cmd.Flags().GetString("team-id")
		name, _ := cmd.Flags().GetString("name")
		colorHex, _ := cmd.Flags().GetString("color")
		stateType, _ := cmd.Flags().GetString("type")

		input := map[string]interface{}{
			"teamId": teamID,
			"name":   name,
			"color":  colorHex,
			"type":   stateType,
		}
		if cmd.Flags().Changed("description") {
			d, _ := cmd.Flags().GetString("description")
			input["description"] = d
		}
		if cmd.Flags().Changed("position") {
			p, _ := cmd.Flags().GetFloat64("position")
			input["position"] = p
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		id, err := client.CreateWorkflowState(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create workflow state: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(map[string]string{"id": id, "name": name})
			return
		}
		output.Success(fmt.Sprintf("Created state %s (%s)", name, id), plaintext, jsonOut)
	},
}

var teamStateUpdateCmd = &cobra.Command{
	Use:   "update STATE-ID",
	Short: "Update a workflow state",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			input["name"] = v
		}
		if cmd.Flags().Changed("color") {
			v, _ := cmd.Flags().GetString("color")
			input["color"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			input["description"] = v
		}
		if cmd.Flags().Changed("position") {
			v, _ := cmd.Flags().GetFloat64("position")
			input["position"] = v
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --color, --description, or --position.", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UpdateWorkflowState(context.Background(), args[0], input); err != nil {
			output.Error(fmt.Sprintf("Failed to update workflow state: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Updated state %s", args[0]), plaintext, jsonOut)
	},
}

var teamStateArchiveCmd = &cobra.Command{
	Use:     "archive STATE-ID",
	Aliases: []string{"rm"},
	Short:   "Archive a workflow state",
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
		if err := client.ArchiveWorkflowState(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to archive workflow state: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Archived state %s", args[0]), plaintext, jsonOut)
	},
}

func init() {
	teamCmd.AddCommand(teamStateCmd)
	teamStateCmd.AddCommand(teamStateCreateCmd)
	teamStateCmd.AddCommand(teamStateUpdateCmd)
	teamStateCmd.AddCommand(teamStateArchiveCmd)

	teamStateCreateCmd.Flags().String("team-id", "", "Team ID (required)")
	teamStateCreateCmd.Flags().String("name", "", "State name (required)")
	teamStateCreateCmd.Flags().String("color", "", "State color (hex, e.g. #e11d48) (required)")
	teamStateCreateCmd.Flags().String("type", "", "State type: backlog|unstarted|started|completed|canceled (required)")
	teamStateCreateCmd.Flags().String("description", "", "State description")
	teamStateCreateCmd.Flags().Float64("position", 0, "State position")
	_ = teamStateCreateCmd.MarkFlagRequired("team-id")
	_ = teamStateCreateCmd.MarkFlagRequired("name")
	_ = teamStateCreateCmd.MarkFlagRequired("color")
	_ = teamStateCreateCmd.MarkFlagRequired("type")

	teamStateUpdateCmd.Flags().String("name", "", "New name")
	teamStateUpdateCmd.Flags().String("color", "", "New color (hex)")
	teamStateUpdateCmd.Flags().String("description", "", "New description")
	teamStateUpdateCmd.Flags().Float64("position", 0, "New position")
}
