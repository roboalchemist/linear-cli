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

// initiativeLabelCmd is the parent command: initiative label
var initiativeLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage initiative labels",
	Long: `Manage Linear initiative labels.

Examples:
  linear-cli initiative label list            # List all initiative labels
  linear-cli initiative label list -l 100     # Limit results`,
}

var initiativeLabelListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List initiative labels",
	Long:    `List all initiative labels in the workspace.`,
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

		labels, err := client.GetInitiativeLabels(context.Background(), nil, limit, "", "", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list initiative labels: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(labels.Nodes)
			return
		}

		if len(labels.Nodes) == 0 {
			if plaintext {
				fmt.Println("No initiative labels found")
			} else {
				fmt.Printf("\n%s No initiative labels found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Initiative Labels")
			fmt.Println("Name\tColor\tDescription")
			for _, l := range labels.Nodes {
				desc := ""
				if l.Description != nil {
					desc = *l.Description
				}
				fmt.Printf("%s\t%s\t%s\n", l.Name, l.Color, desc)
			}
		} else {
			headers := []string{"Name", "Color", "Description"}
			rows := [][]string{}

			for _, l := range labels.Nodes {
				desc := ""
				if l.Description != nil {
					desc = *l.Description
				}
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(l.Name),
					l.Color,
					desc,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d initiative labels\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(labels.Nodes))
		}
	},
}

func init() {
	initiativeCmd.AddCommand(initiativeLabelCmd)
	initiativeLabelCmd.AddCommand(initiativeLabelListCmd)

	initiativeLabelListCmd.Flags().IntP("limit", "l", 50, "Maximum number of initiative labels to return")
}
