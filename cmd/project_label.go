package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// projectLabelCmd is the parent command: project label
var projectLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage project labels",
	Long: `Manage Linear project labels.

Examples:
  linear-cli project label list            # List all project labels
  linear-cli project label list -l 100     # Limit results`,
}

var projectLabelListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List project labels",
	Long:    `List all project labels in the workspace.`,
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

		labels, err := client.GetProjectLabels(context.Background(), nil, limit, "", "", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list project labels: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(labels.Nodes)
			return
		}

		if len(labels.Nodes) == 0 {
			if plaintext {
				fmt.Println("No project labels found")
			} else {
				fmt.Printf("\n%s No project labels found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Project Labels")
			fmt.Println("Name\tColor\tTeam\tDescription")
			for _, l := range labels.Nodes {
				desc := ""
				if l.Description != nil {
					desc = *l.Description
				}
				team := ""
				if l.Team != nil {
					team = l.Team.Key
				}
				fmt.Printf("%s\t%s\t%s\t%s\n", l.Name, l.Color, team, desc)
			}
		} else {
			headers := []string{"Name", "Color", "Team", "Description"}
			rows := [][]string{}

			for _, l := range labels.Nodes {
				desc := ""
				if l.Description != nil {
					desc = *l.Description
				}
				team := color.New(color.FgWhite, color.Faint).Sprint("workspace")
				if l.Team != nil {
					team = l.Team.Key
				}
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(l.Name),
					l.Color,
					team,
					desc,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d project labels\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(labels.Nodes))
		}
	},
}

func init() {
	projectCmd.AddCommand(projectLabelCmd)
	projectLabelCmd.AddCommand(projectLabelListCmd)

	projectLabelListCmd.Flags().IntP("limit", "l", 50, "Maximum number of project labels to return")
}
