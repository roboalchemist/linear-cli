package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage Linear webhooks",
	Long: `Manage Linear webhooks.

Examples:
  linear-cli webhook list            # List all webhooks
  linear-cli webhook get WEBHOOK-ID  # Get webhook details`,
}

var webhookListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List webhooks",
	Long:    `List all webhooks for the current workspace.`,
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

		webhooks, err := client.GetWebhooks(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list webhooks: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(webhooks.Nodes)
			return
		}

		if len(webhooks.Nodes) == 0 {
			if plaintext {
				fmt.Println("No webhooks found")
			} else {
				fmt.Printf("\n%s No webhooks found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"URL", "Enabled", "Team", "Resource Types", "Created"}
		rows := [][]string{}
		for _, w := range webhooks.Nodes {
			enabled := "no"
			if w.Enabled {
				enabled = "yes"
			}
			rows = append(rows, []string{
				webhookURL(w),
				enabled,
				webhookTeam(w),
				strings.Join(w.ResourceTypes, ", "),
				formatDateShort(w.CreatedAt.Format(time.RFC3339)),
			})
		}

		if plaintext {
			fmt.Println("# Webhooks")
			fmt.Println(strings.Join(headers, "\t"))
			for _, r := range rows {
				fmt.Println(strings.Join(r, "\t"))
			}
			return
		}

		output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
	},
}

var webhookGetCmd = &cobra.Command{
	Use:     "get WEBHOOK-ID",
	Aliases: []string{"show"},
	Short:   "Get webhook details",
	Long:    `Get details for a specific webhook.`,
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
		webhook, err := client.GetWebhook(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get webhook: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(webhook)
			return
		}

		enabled := "no"
		if webhook.Enabled {
			enabled = "yes"
		}

		if plaintext {
			fmt.Printf("# Webhook %s\n", webhook.ID)
			fmt.Printf("URL: %s\n", webhookURL(*webhook))
			fmt.Printf("Enabled: %s\n", enabled)
			fmt.Printf("Team: %s\n", webhookTeam(*webhook))
			fmt.Printf("Resource Types: %s\n", strings.Join(webhook.ResourceTypes, ", "))
			fmt.Printf("All Public Teams: %t\n", webhook.AllPublicTeams)
			fmt.Printf("Created: %s\n", formatDateShort(webhook.CreatedAt.Format(time.RFC3339)))
			return
		}

		fmt.Printf("\n%s Webhook %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("🔗"),
			color.New(color.FgWhite, color.Bold).Sprint(webhook.ID))
		fmt.Printf("   URL: %s\n", webhookURL(*webhook))
		if webhook.Enabled {
			fmt.Printf("   Enabled: %s\n", color.New(color.FgGreen).Sprint(enabled))
		} else {
			fmt.Printf("   Enabled: %s\n", color.New(color.FgRed).Sprint(enabled))
		}
		fmt.Printf("   Team: %s\n", webhookTeam(*webhook))
		fmt.Printf("   Resource Types: %s\n", strings.Join(webhook.ResourceTypes, ", "))
		fmt.Printf("   Created: %s\n", formatDateShort(webhook.CreatedAt.Format(time.RFC3339)))
	},
}

// webhookURL returns the webhook destination URL, falling back to its label.
func webhookURL(w api.Webhook) string {
	if w.URL != nil && *w.URL != "" {
		return *w.URL
	}
	if w.Label != nil && *w.Label != "" {
		return *w.Label
	}
	return "-"
}

// webhookTeam returns a human-readable description of the webhook's scope.
func webhookTeam(w api.Webhook) string {
	if w.Team != nil {
		return w.Team.Key
	}
	if w.AllPublicTeams {
		return "all public teams"
	}
	if len(w.TeamIds) > 0 {
		return fmt.Sprintf("%d teams", len(w.TeamIds))
	}
	return "workspace"
}

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(webhookListCmd)
	webhookCmd.AddCommand(webhookGetCmd)

	webhookListCmd.Flags().IntP("limit", "l", 50, "Maximum number of webhooks to return")
}
