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

var externalCmd = &cobra.Command{
	Use:   "external",
	Short: "Manage external users",
	Long: `Inspect external users — people who interact with Linear through
integrated services (Slack, Jira, GitHub, ...) without a Linear account.

Examples:
  linear-cli external user list
  linear-cli external user list --limit 10`,
}

var externalUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Manage external users",
	Long:    `List external users for the workspace.`,
}

var externalUserListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List external users",
	Long:    `List all external users in the workspace.`,
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

		users, err := client.GetExternalUsers(context.Background(), limit, "", "createdAt")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list external users: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(users.Nodes)
			return
		}

		if len(users.Nodes) == 0 {
			output.Info("No external users found", plaintext, jsonOut)
			return
		}

		// NOTE: Linear's ExternalUser schema type does not expose a provider /
		// integration-service field, so the Provider column is rendered as "-".
		headers := []string{"Name", "Email", "Provider", "Created"}
		rows := make([][]string, len(users.Nodes))
		for i, u := range users.Nodes {
			email := "-"
			if u.Email != nil && *u.Email != "" {
				email = *u.Email
			}
			rows[i] = []string{
				truncateString(u.Name, 30),
				email,
				"-",
				u.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		if users.PageInfo.HasNextPage && !plaintext {
			fmt.Printf("\n%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

func init() {
	rootCmd.AddCommand(externalCmd)
	externalCmd.AddCommand(externalUserCmd)
	externalUserCmd.AddCommand(externalUserListCmd)

	externalUserListCmd.Flags().IntP("limit", "l", 50, "Maximum number of external users to return")
}
