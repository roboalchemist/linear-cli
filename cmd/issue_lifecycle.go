package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ---------------------------------------------------------------------------
// issue unarchive / delete
// ---------------------------------------------------------------------------

var issueUnarchiveCmd = &cobra.Command{
	Use:   "unarchive ISSUE-ID",
	Short: "Restore an archived issue",
	Long:  `Unarchive an issue, returning it to the active issue list.`,
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
		if err := client.UnarchiveIssue(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to unarchive issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Unarchived issue %s", args[0]), plaintext, jsonOut)
	},
}

var issueDeleteCmd = &cobra.Command{
	Use:     "delete ISSUE-ID",
	Aliases: []string{"rm"},
	Short:   "Delete an issue",
	Long: `Delete an issue. By default the issue is trashed (recoverable); pass
--permanent to delete it permanently.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		permanent, _ := cmd.Flags().GetBool("permanent")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.DeleteIssue(context.Background(), args[0], permanent); err != nil {
			output.Error(fmt.Sprintf("Failed to delete issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if permanent {
			output.Success(fmt.Sprintf("Permanently deleted issue %s", args[0]), plaintext, jsonOut)
		} else {
			output.Success(fmt.Sprintf("Trashed issue %s", args[0]), plaintext, jsonOut)
		}
	},
}

// ---------------------------------------------------------------------------
// issue label add / remove
// ---------------------------------------------------------------------------

var issueLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Attach or detach labels on an issue",
	Long:  `Attach or detach labels on an existing issue.`,
}

var issueLabelAddCmd = &cobra.Command{
	Use:   "add ISSUE-ID LABEL-ID [LABEL-ID...]",
	Short: "Add label(s) to an issue",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		issueID := args[0]
		added := 0
		for _, labelID := range args[1:] {
			if err := client.AddIssueLabel(context.Background(), issueID, labelID); err != nil {
				output.Error(fmt.Sprintf("Failed to add label %s: %v", labelID, err), plaintext, jsonOut)
				os.Exit(1)
			}
			added++
		}
		output.Success(fmt.Sprintf("Added %d label(s) to %s", added, issueID), plaintext, jsonOut)
	},
}

var issueLabelRemoveCmd = &cobra.Command{
	Use:     "remove ISSUE-ID LABEL-ID [LABEL-ID...]",
	Aliases: []string{"rm"},
	Short:   "Remove label(s) from an issue",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		issueID := args[0]
		removed := 0
		for _, labelID := range args[1:] {
			if err := client.RemoveIssueLabel(context.Background(), issueID, labelID); err != nil {
				output.Error(fmt.Sprintf("Failed to remove label %s: %v", labelID, err), plaintext, jsonOut)
				os.Exit(1)
			}
			removed++
		}
		output.Success(fmt.Sprintf("Removed %d label(s) from %s", removed, issueID), plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// issue subscribe / unsubscribe / share / unshare
// ---------------------------------------------------------------------------

var issueSubscribeCmd = &cobra.Command{
	Use:   "subscribe ISSUE-ID",
	Short: "Subscribe to an issue",
	Long:  `Subscribe to an issue's notifications. Defaults to the current user.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		email, _ := cmd.Flags().GetString("email")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.SubscribeIssue(context.Background(), args[0], email); err != nil {
			output.Error(fmt.Sprintf("Failed to subscribe: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Subscribed to %s", args[0]), plaintext, jsonOut)
	},
}

var issueUnsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe ISSUE-ID",
	Short: "Unsubscribe from an issue",
	Long:  `Unsubscribe from an issue's notifications. Defaults to the current user.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		email, _ := cmd.Flags().GetString("email")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UnsubscribeIssue(context.Background(), args[0], email); err != nil {
			output.Error(fmt.Sprintf("Failed to unsubscribe: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Unsubscribed from %s", args[0]), plaintext, jsonOut)
	},
}

var issueShareCmd = &cobra.Command{
	Use:   "share ISSUE-ID --user USER-ID",
	Short: "Share an issue with a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		userID, _ := cmd.Flags().GetString("user")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.ShareIssue(context.Background(), args[0], userID); err != nil {
			output.Error(fmt.Sprintf("Failed to share issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Shared %s with %s", args[0], userID), plaintext, jsonOut)
	},
}

var issueUnshareCmd = &cobra.Command{
	Use:   "unshare ISSUE-ID --user USER-ID",
	Short: "Remove a user's access to a shared issue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		userID, _ := cmd.Flags().GetString("user")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UnshareIssue(context.Background(), args[0], userID); err != nil {
			output.Error(fmt.Sprintf("Failed to unshare issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Removed %s's access to %s", userID, args[0]), plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// issue reminder add / remove
// ---------------------------------------------------------------------------

var issueReminderCmd = &cobra.Command{
	Use:   "reminder",
	Short: "Manage issue reminders",
	Long:  `Set or clear a reminder for the current user on an issue.`,
}

var issueReminderAddCmd = &cobra.Command{
	Use:   "add ISSUE-ID --at RFC3339",
	Short: "Set a reminder on an issue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		at, _ := cmd.Flags().GetString("at")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.AddIssueReminder(context.Background(), args[0], at); err != nil {
			output.Error(fmt.Sprintf("Failed to add reminder: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Reminder set on %s for %s", args[0], at), plaintext, jsonOut)
	},
}

var issueReminderRemoveCmd = &cobra.Command{
	Use:     "remove ISSUE-ID",
	Aliases: []string{"rm"},
	Short:   "Clear a reminder on an issue",
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
		if err := client.RemoveIssueReminder(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to remove reminder: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Removed reminder on %s", args[0]), plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// issue batch create / update
// ---------------------------------------------------------------------------

var issueBatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Bulk create or update issues",
	Long: `Bulk issue operations.

Examples:
  linear-cli issue batch create --file issues.json   # JSON array of IssueCreateInput
  linear-cli issue batch update ID1 ID2 --state-id STATE_ID --priority 2`,
}

var issueBatchCreateCmd = &cobra.Command{
	Use:   "create --file FILE",
	Short: "Create many issues from a JSON array",
	Long: `Create many issues at once from a JSON file containing an array of
IssueCreateInput objects (use "-" to read stdin).`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		file, _ := cmd.Flags().GetString("file")

		raw, err := readContentFromFile(file)
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}
		var issues []map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &issues); err != nil {
			output.Error(fmt.Sprintf("Invalid JSON (expected an array of IssueCreateInput): %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if len(issues) == 0 {
			output.Error("No issues provided", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.BatchCreateIssues(context.Background(), issues); err != nil {
			output.Error(fmt.Sprintf("Failed to batch create issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Created %d issue(s)", len(issues)), plaintext, jsonOut)
	},
}

var issueBatchUpdateCmd = &cobra.Command{
	Use:   "update ISSUE-ID [ISSUE-ID...]",
	Short: "Apply the same update to many issues",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		input := map[string]interface{}{}
		if cmd.Flags().Changed("state-id") {
			v, _ := cmd.Flags().GetString("state-id")
			input["stateId"] = v
		}
		if cmd.Flags().Changed("assignee-id") {
			v, _ := cmd.Flags().GetString("assignee-id")
			input["assigneeId"] = v
		}
		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetInt("priority")
			input["priority"] = v
		}
		if cmd.Flags().Changed("project-id") {
			v, _ := cmd.Flags().GetString("project-id")
			input["projectId"] = v
		}
		if cmd.Flags().Changed("cycle-id") {
			v, _ := cmd.Flags().GetString("cycle-id")
			input["cycleId"] = v
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --state-id, --assignee-id, --priority, --project-id, or --cycle-id.",
				plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.BatchUpdateIssues(context.Background(), args, input); err != nil {
			output.Error(fmt.Sprintf("Failed to batch update issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Updated %d issue(s)", len(args)), plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// issue comment resolve / unresolve
// ---------------------------------------------------------------------------

var commentResolveCmd = &cobra.Command{
	Use:   "resolve COMMENT-ID",
	Short: "Resolve a comment thread",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		resolving, _ := cmd.Flags().GetString("resolving-comment")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.ResolveComment(context.Background(), args[0], resolving); err != nil {
			output.Error(fmt.Sprintf("Failed to resolve comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Resolved comment", plaintext, jsonOut)
	},
}

var commentUnresolveCmd = &cobra.Command{
	Use:   "unresolve COMMENT-ID",
	Short: "Reopen a resolved comment thread",
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
		if err := client.UnresolveComment(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to unresolve comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Reopened comment", plaintext, jsonOut)
	},
}

func init() {
	// issue lifecycle
	issueCmd.AddCommand(issueUnarchiveCmd)
	issueCmd.AddCommand(issueDeleteCmd)
	issueCmd.AddCommand(issueLabelCmd)
	issueLabelCmd.AddCommand(issueLabelAddCmd)
	issueLabelCmd.AddCommand(issueLabelRemoveCmd)
	issueCmd.AddCommand(issueSubscribeCmd)
	issueCmd.AddCommand(issueUnsubscribeCmd)
	issueCmd.AddCommand(issueShareCmd)
	issueCmd.AddCommand(issueUnshareCmd)
	issueCmd.AddCommand(issueReminderCmd)
	issueReminderCmd.AddCommand(issueReminderAddCmd)
	issueReminderCmd.AddCommand(issueReminderRemoveCmd)
	issueCmd.AddCommand(issueBatchCmd)
	issueBatchCmd.AddCommand(issueBatchCreateCmd)
	issueBatchCmd.AddCommand(issueBatchUpdateCmd)

	issueDeleteCmd.Flags().Bool("permanent", false, "Permanently delete (skip trash)")
	issueSubscribeCmd.Flags().String("email", "", "Email of the user to subscribe (default: you)")
	issueUnsubscribeCmd.Flags().String("email", "", "Email of the user to unsubscribe (default: you)")
	issueShareCmd.Flags().String("user", "", "User ID to share with (required)")
	issueUnshareCmd.Flags().String("user", "", "User ID to remove (required)")
	_ = issueShareCmd.MarkFlagRequired("user")
	_ = issueUnshareCmd.MarkFlagRequired("user")
	issueReminderAddCmd.Flags().String("at", "", "Reminder time (RFC3339, e.g. 2026-02-01T09:00:00Z) (required)")
	_ = issueReminderAddCmd.MarkFlagRequired("at")
	issueBatchCreateCmd.Flags().String("file", "", "Path to JSON array of IssueCreateInput (use - for stdin) (required)")
	_ = issueBatchCreateCmd.MarkFlagRequired("file")
	issueBatchUpdateCmd.Flags().String("state-id", "", "New workflow state ID")
	issueBatchUpdateCmd.Flags().String("assignee-id", "", "New assignee user ID")
	issueBatchUpdateCmd.Flags().Int("priority", 0, "New priority 0-4")
	issueBatchUpdateCmd.Flags().String("project-id", "", "New project ID")
	issueBatchUpdateCmd.Flags().String("cycle-id", "", "New cycle ID")

	// issue comment lifecycle
	commentCmd.AddCommand(commentResolveCmd)
	commentCmd.AddCommand(commentUnresolveCmd)
	commentResolveCmd.Flags().String("resolving-comment", "", "ID of the child comment that resolves the thread")
}
