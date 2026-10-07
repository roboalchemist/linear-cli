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

var webhookCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a webhook",
	Long: `Create a new webhook subscription.

--url and at least one --resource-types value are required. Use
--all-public-teams for a workspace-wide webhook, or --team-id for one team.

Examples:
  linear-cli webhook create --url https://example.com/hook --resource-types Issue,Comment
  linear-cli webhook create --url https://example.com/hook --all-public-teams --resource-types Issue`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		url, _ := cmd.Flags().GetString("url")
		teamID, _ := cmd.Flags().GetString("team-id")
		resourceTypes, _ := cmd.Flags().GetStringSlice("resource-types")
		allPublicTeams, _ := cmd.Flags().GetBool("all-public-teams")
		enabled, _ := cmd.Flags().GetBool("enabled")
		label, _ := cmd.Flags().GetString("label")
		secret, _ := cmd.Flags().GetString("secret")
		inputJSON, _ := cmd.Flags().GetString("input-json")

		input := map[string]interface{}{"url": url}
		if teamID != "" {
			input["teamId"] = teamID
		}
		if len(resourceTypes) > 0 {
			input["resourceTypes"] = resourceTypes
		}
		if allPublicTeams {
			input["allPublicTeams"] = allPublicTeams
		}
		if cmd.Flags().Changed("enabled") {
			input["enabled"] = enabled
		}
		if label != "" {
			input["label"] = label
		}
		if secret != "" {
			input["secret"] = secret
		}
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if _, ok := input["resourceTypes"]; !ok {
			output.Error("--resource-types is required (or provide resourceTypes via --input-json).", plaintext, jsonOut)
			os.Exit(1)
		}

		webhook, err := client.CreateWebhook(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create webhook: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(webhook)
		} else {
			output.Success(fmt.Sprintf("Created webhook %s",
				color.New(color.FgWhite, color.Bold).Sprint(webhookURL(*webhook))), plaintext, jsonOut)
		}
	},
}

var webhookUpdateCmd = &cobra.Command{
	Use:     "update WEBHOOK-ID",
	Aliases: []string{"edit"},
	Short:   "Update a webhook",
	Long: `Update a webhook's URL, resources, or enabled state.

Use --enabled or --disabled to toggle the webhook.

Examples:
  linear-cli webhook update WEBHOOK-ID --url https://example.com/hook2
  linear-cli webhook update WEBHOOK-ID --disabled`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		input := map[string]interface{}{}
		if cmd.Flags().Changed("url") {
			v, _ := cmd.Flags().GetString("url")
			input["url"] = v
		}
		if cmd.Flags().Changed("label") {
			v, _ := cmd.Flags().GetString("label")
			input["label"] = v
		}
		if cmd.Flags().Changed("secret") {
			v, _ := cmd.Flags().GetString("secret")
			input["secret"] = v
		}
		if cmd.Flags().Changed("resource-types") {
			v, _ := cmd.Flags().GetStringSlice("resource-types")
			input["resourceTypes"] = v
		}
		if cmd.Flags().Changed("enabled") {
			v, _ := cmd.Flags().GetBool("enabled")
			input["enabled"] = v
		}
		if cmd.Flags().Changed("disabled") {
			v, _ := cmd.Flags().GetBool("disabled")
			input["enabled"] = !v
		}
		inputJSON, _ := cmd.Flags().GetString("input-json")
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --url, --label, --secret, --resource-types, --enabled, --disabled, or --input-json.", plaintext, jsonOut)
			os.Exit(1)
		}

		webhook, err := client.UpdateWebhook(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update webhook: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(webhook)
		} else {
			output.Success(fmt.Sprintf("Updated webhook %s",
				color.New(color.FgWhite, color.Bold).Sprint(webhookURL(*webhook))), plaintext, jsonOut)
		}
	},
}

var webhookDeleteCmd = &cobra.Command{
	Use:     "delete WEBHOOK-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a webhook",
	Long:    `Delete a webhook.`,
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
		if err := client.DeleteWebhook(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete webhook: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted webhook", plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(webhookListCmd)
	webhookCmd.AddCommand(webhookGetCmd)
	webhookCmd.AddCommand(webhookCreateCmd)
	webhookCmd.AddCommand(webhookUpdateCmd)
	webhookCmd.AddCommand(webhookDeleteCmd)

	webhookListCmd.Flags().IntP("limit", "l", 50, "Maximum number of webhooks to return")

	// Create flags
	webhookCreateCmd.Flags().String("url", "", "Webhook destination URL (required)")
	webhookCreateCmd.Flags().String("team-id", "", "Team ID to scope the webhook to")
	webhookCreateCmd.Flags().StringSlice("resource-types", nil, "Resource types to subscribe to (required, e.g. Issue,Comment)")
	webhookCreateCmd.Flags().Bool("all-public-teams", false, "Enable the webhook for all public teams")
	webhookCreateCmd.Flags().Bool("enabled", true, "Whether the webhook is enabled")
	webhookCreateCmd.Flags().String("label", "", "Label for the webhook")
	webhookCreateCmd.Flags().String("secret", "", "Secret token used to sign the webhook payload")
	webhookCreateCmd.Flags().String("input-json", "", "Additional WebhookCreateInput fields as a JSON object")
	_ = webhookCreateCmd.MarkFlagRequired("url")

	// Update flags
	webhookUpdateCmd.Flags().String("url", "", "New webhook destination URL")
	webhookUpdateCmd.Flags().String("label", "", "New label")
	webhookUpdateCmd.Flags().String("secret", "", "New secret token")
	webhookUpdateCmd.Flags().StringSlice("resource-types", nil, "New resource types to subscribe to")
	webhookUpdateCmd.Flags().Bool("enabled", true, "Enable the webhook")
	webhookUpdateCmd.Flags().Bool("disabled", false, "Disable the webhook")
	webhookUpdateCmd.Flags().String("input-json", "", "Additional WebhookUpdateInput fields as a JSON object")
}
