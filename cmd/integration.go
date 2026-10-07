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

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Manage Linear workspace integrations",
	Long: `List Linear workspace integrations (Slack, GitHub, Jira, Figma, ...).

Examples:
  linear-cli integration list
  linear-cli integration list --limit 10`,
}

var integrationListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List integrations",
	Long:    `List all integrations configured for the workspace.`,
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

		integrations, err := client.GetIntegrations(context.Background(), limit, "", "createdAt")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list integrations: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(integrations.Nodes)
			return
		}

		if len(integrations.Nodes) == 0 {
			output.Info("No integrations found", plaintext, jsonOut)
			return
		}

		headers := []string{"Service", "Team", "Created"}
		rows := make([][]string, len(integrations.Nodes))
		for i, ig := range integrations.Nodes {
			team := "Workspace"
			if ig.Team != nil {
				team = ig.Team.Key
			}
			rows[i] = []string{
				color.New(color.FgCyan).Sprint(ig.Service),
				team,
				ig.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		if integrations.PageInfo.HasNextPage && !plaintext {
			fmt.Printf("\n%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

func init() {
	rootCmd.AddCommand(integrationCmd)
	integrationCmd.AddCommand(integrationListCmd)

	integrationListCmd.Flags().IntP("limit", "l", 50, "Maximum number of integrations to return")
}
