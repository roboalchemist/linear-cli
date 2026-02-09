package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dorkitude/linear-cli/pkg/api"
	"github.com/dorkitude/linear-cli/pkg/auth"
	"github.com/dorkitude/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage Linear labels",
	Long: `Manage issue labels including listing and creating labels.

Examples:
  linear-cli label list                        # List all labels
  linear-cli label list --team ROB             # List labels for a team
  linear-cli label create --name "bug" --color "#e11d48" --team TEAM-ID`,
}

var labelListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List labels",
	Long:    `List issue labels, optionally filtered by team.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")
		teamKey, _ := cmd.Flags().GetString("team")

		filter := map[string]interface{}{}
		if teamKey != "" {
			filter["team"] = map[string]interface{}{
				"key": map[string]interface{}{"eq": teamKey},
			}
		}

		labels, err := client.GetLabels(context.Background(), filter, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list labels: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(labels.Nodes)
			return
		}

		if len(labels.Nodes) == 0 {
			if plaintext {
				fmt.Println("No labels found")
			} else {
				fmt.Printf("\n%s No labels found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Labels")
			fmt.Println("Name\tColor\tDescription\tParent")
			for _, l := range labels.Nodes {
				desc := ""
				if l.Description != nil {
					desc = *l.Description
				}
				parent := ""
				if l.Parent != nil {
					parent = l.Parent.Name
				}
				fmt.Printf("%s\t%s\t%s\t%s\n", l.Name, l.Color, desc, parent)
			}
		} else {
			headers := []string{"Name", "Color", "Description", "Parent"}
			rows := [][]string{}

			for _, l := range labels.Nodes {
				desc := ""
				if l.Description != nil {
					desc = *l.Description
				}
				parent := ""
				if l.Parent != nil {
					parent = l.Parent.Name
				}

				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(l.Name),
					l.Color,
					desc,
					parent,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d labels\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(labels.Nodes))
		}
	},
}

var labelCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new label",
	Long:    `Create a new issue label.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		name, _ := cmd.Flags().GetString("name")
		labelColor, _ := cmd.Flags().GetString("color")
		teamID, _ := cmd.Flags().GetString("team-id")

		input := map[string]interface{}{
			"name": name,
		}
		if labelColor != "" {
			input["color"] = labelColor
		}
		if teamID != "" {
			input["teamId"] = teamID
		}

		label, err := client.CreateLabel(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create label: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(label)
		} else if plaintext {
			fmt.Printf("Created label: %s (%s)\n", label.Name, label.Color)
		} else {
			output.Success(fmt.Sprintf("Created label %s",
				color.New(color.FgWhite, color.Bold).Sprint(label.Name)), plaintext, jsonOut)
		}
	},
}

var labelUpdateCmd = &cobra.Command{
	Use:     "update LABEL-ID",
	Aliases: []string{"edit"},
	Short:   "Update a label",
	Long:    `Update a label's name, color, or description.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		labelID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			input["name"] = name
		}
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}
		if cmd.Flags().Changed("description") {
			d, _ := cmd.Flags().GetString("description")
			input["description"] = d
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --color, or --description.", plaintext, jsonOut)
			os.Exit(1)
		}

		label, err := client.UpdateLabel(context.Background(), labelID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update label: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(label)
		} else {
			output.Success(fmt.Sprintf("Updated label %s",
				color.New(color.FgWhite, color.Bold).Sprint(label.Name)), plaintext, jsonOut)
		}
	},
}

var labelDeleteCmd = &cobra.Command{
	Use:     "delete LABEL-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a label",
	Long:    `Delete a label.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		labelID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		err = client.DeleteLabel(context.Background(), labelID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete label: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted label", plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(labelListCmd)
	labelCmd.AddCommand(labelCreateCmd)
	labelCmd.AddCommand(labelUpdateCmd)
	labelCmd.AddCommand(labelDeleteCmd)

	// List flags
	labelListCmd.Flags().IntP("limit", "l", 50, "Maximum number of labels to return")
	labelListCmd.Flags().StringP("team", "t", "", "Filter by team key")

	// Create flags
	labelCreateCmd.Flags().StringP("name", "n", "", "Label name (required)")
	labelCreateCmd.Flags().StringP("color", "c", "", "Label color (hex, e.g., #e11d48)")
	labelCreateCmd.Flags().String("team-id", "", "Team ID to scope the label to")
	_ = labelCreateCmd.MarkFlagRequired("name")

	// Update flags
	labelUpdateCmd.Flags().StringP("name", "n", "", "New label name")
	labelUpdateCmd.Flags().StringP("color", "c", "", "New label color (hex)")
	labelUpdateCmd.Flags().StringP("description", "d", "", "New label description")
}
