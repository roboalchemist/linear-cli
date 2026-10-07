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

var subscriptionCmd = &cobra.Command{
	Use:     "subscription",
	Aliases: []string{"subscriptions", "notif-subscription"},
	Short:   "Manage notification subscriptions",
	Long: `Inspect the authenticated user's notification subscriptions for teams,
projects, cycles, initiatives, labels, customers, and users.

Examples:
  linear-cli subscription list
  linear-cli subscription list --limit 10`,
}

var subscriptionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List notification subscriptions",
	Long:    `List the authenticated user's notification subscriptions.`,
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

		subs, err := client.GetNotificationSubscriptions(context.Background(), limit, "", "createdAt", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list notification subscriptions: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(subs.Nodes)
			return
		}

		if len(subs.Nodes) == 0 {
			output.Info("No notification subscriptions found", plaintext, jsonOut)
			return
		}

		headers := []string{"Target", "Type", "Active", "Created"}
		rows := make([][]string, len(subs.Nodes))
		for i, s := range subs.Nodes {
			target, kind := subscriptionTarget(s)
			active := color.New(color.FgRed).Sprint("no")
			if s.Active {
				active = color.New(color.FgGreen).Sprint("yes")
			}
			rows[i] = []string{
				truncateString(target, 40),
				kind,
				active,
				s.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		if subs.PageInfo.HasNextPage && !plaintext {
			fmt.Printf("\n%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

// subscriptionTarget returns a human-readable target name and its entity type.
func subscriptionTarget(s api.NotificationSubscription) (string, string) {
	switch {
	case s.Team != nil:
		return s.Team.Name, "team"
	case s.Project != nil:
		return s.Project.Name, "project"
	case s.Cycle != nil:
		return s.Cycle.Name, "cycle"
	case s.Initiative != nil:
		return s.Initiative.Name, "initiative"
	case s.Label != nil:
		return s.Label.Name, "label"
	case s.User != nil:
		return s.User.Name, "user"
	default:
		return s.ID, "unknown"
	}
}

func init() {
	rootCmd.AddCommand(subscriptionCmd)
	subscriptionCmd.AddCommand(subscriptionListCmd)

	subscriptionListCmd.Flags().IntP("limit", "l", 50, "Maximum number of subscriptions to return")
}
