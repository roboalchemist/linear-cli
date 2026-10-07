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

// ---------------------------------------------------------------------------
// oauth app list
// ---------------------------------------------------------------------------

var oauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "Manage OAuth applications",
	Long:  `Inspect OAuth applications configured for the workspace.`,
}

var oauthAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage OAuth applications",
}

var oauthAppListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List OAuth applications",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		apps, err := client.GetOAuthApplications(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list OAuth applications: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(apps)
			return
		}
		if len(apps) == 0 {
			output.Info("No OAuth applications found", plaintext, jsonOut)
			return
		}
		rows := [][]string{}
		for _, a := range apps {
			rows = append(rows, []string{a.Name, a.ClientID, a.Distribution, a.CreatedAt})
		}
		output.Table(output.TableData{
			Headers: []string{"Name", "Client ID", "Distribution", "Created"},
			Rows:    rows,
		}, plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// organization invite list
// ---------------------------------------------------------------------------

var organizationInviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Manage organization invites",
}

var organizationInviteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List organization invites",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		limit, _ := cmd.Flags().GetInt("limit")
		includeArchived, _ := cmd.Flags().GetBool("include-archived")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		invites, err := client.GetOrganizationInvites(context.Background(), limit, "", includeArchived)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list invites: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(invites.Nodes)
			return
		}
		if len(invites.Nodes) == 0 {
			output.Info("No organization invites found", plaintext, jsonOut)
			return
		}
		rows := [][]string{}
		for _, i := range invites.Nodes {
			rows = append(rows, []string{i.Email, i.Role, fmt.Sprintf("%t", i.External), i.CreatedAt, i.AcceptedAt})
		}
		output.Table(output.TableData{
			Headers: []string{"Email", "Role", "External", "Created", "Accepted"},
			Rows:    rows,
		}, plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// audit types
// ---------------------------------------------------------------------------

var auditTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List audit log entry types",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		types, err := client.GetAuditEntryTypes(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list audit types: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(types)
			return
		}
		rows := [][]string{}
		for _, t := range types {
			rows = append(rows, []string{t.Type, t.Description})
		}
		output.Table(output.TableData{
			Headers: []string{"Type", "Description"},
			Rows:    rows,
		}, plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// team membership list
// ---------------------------------------------------------------------------

var teamMembershipCmd = &cobra.Command{
	Use:   "membership",
	Short: "Manage team membership records",
}

var teamMembershipListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List team memberships",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		limit, _ := cmd.Flags().GetInt("limit")
		includeArchived, _ := cmd.Flags().GetBool("include-archived")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		ms, err := client.GetTeamMemberships(context.Background(), limit, "", includeArchived)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list team memberships: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(ms.Nodes)
			return
		}
		if len(ms.Nodes) == 0 {
			output.Info("No team memberships found", plaintext, jsonOut)
			return
		}
		rows := [][]string{}
		for _, m := range ms.Nodes {
			teamKey, userName, userEmail := "", "", ""
			if m.Team != nil {
				teamKey = m.Team.Key
			}
			if m.User != nil {
				userName, userEmail = m.User.Name, m.User.Email
			}
			owner := ""
			if m.Owner {
				owner = "yes"
			}
			rows = append(rows, []string{teamKey, userName, userEmail, owner})
		}
		output.Table(output.TableData{
			Headers: []string{"Team", "User", "Email", "Owner"},
			Rows:    rows,
		}, plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// Emoji create / delete + webhook rotate-secret
// ---------------------------------------------------------------------------

var emojiCreateCmd = &cobra.Command{
	Use:   "create --name NAME [--url URL]",
	Short: "Create a custom emoji",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.CreateEmoji(context.Background(), name, url); err != nil {
			output.Error(fmt.Sprintf("Failed to create emoji: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Created emoji :%s:", name), plaintext, jsonOut)
	},
}

var emojiDeleteCmd = &cobra.Command{
	Use:     "delete EMOJI-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a custom emoji",
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
		if err := client.DeleteEmoji(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete emoji: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Deleted emoji", plaintext, jsonOut)
	},
}

var webhookRotateSecretCmd = &cobra.Command{
	Use:   "rotate-secret WEBHOOK-ID",
	Short: "Rotate a webhook's signing secret",
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
		wh, err := client.RotateWebhookSecret(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to rotate webhook secret: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(wh)
			return
		}
		output.Success("Rotated webhook secret", plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(oauthCmd)
	oauthCmd.AddCommand(oauthAppCmd)
	oauthAppCmd.AddCommand(oauthAppListCmd)

	organizationCmd.AddCommand(organizationInviteCmd)
	organizationInviteCmd.AddCommand(organizationInviteListCmd)
	organizationInviteListCmd.Flags().IntP("limit", "l", 50, "Maximum number of invites to return")
	organizationInviteListCmd.Flags().Bool("include-archived", false, "Include archived invites")

	auditCmd.AddCommand(auditTypesCmd)

	teamCmd.AddCommand(teamMembershipCmd)
	teamMembershipCmd.AddCommand(teamMembershipListCmd)
	teamMembershipListCmd.Flags().IntP("limit", "l", 50, "Maximum number of memberships to return")
	teamMembershipListCmd.Flags().Bool("include-archived", false, "Include archived memberships")

	emojiCmd.AddCommand(emojiCreateCmd)
	emojiCmd.AddCommand(emojiDeleteCmd)
	emojiCreateCmd.Flags().String("name", "", "Emoji name (without colons) (required)")
	emojiCreateCmd.Flags().String("url", "", "Image URL for the emoji")
	_ = emojiCreateCmd.MarkFlagRequired("name")

	webhookCmd.AddCommand(webhookRotateSecretCmd)
}
