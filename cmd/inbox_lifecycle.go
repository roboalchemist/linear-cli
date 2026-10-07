package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// notificationEntityInput builds a NotificationEntityInput from flags. An empty
// map means "apply to everything".
func notificationEntityInput(cmd *cobra.Command) map[string]interface{} {
	input := map[string]interface{}{}
	for _, f := range []struct{ flag, field string }{
		{"issue", "issueId"},
		{"project", "projectId"},
		{"initiative", "initiativeId"},
		{"project-update", "projectUpdateId"},
		{"initiative-update", "initiativeUpdateId"},
	} {
		if v, _ := cmd.Flags().GetString(f.flag); v != "" {
			input[f.field] = v
		}
	}
	return input
}

func addNotificationEntityFlags(cmd *cobra.Command) {
	cmd.Flags().String("issue", "", "Scope to a single issue ID")
	cmd.Flags().String("project", "", "Scope to a single project ID")
	cmd.Flags().String("initiative", "", "Scope to a single initiative ID")
	cmd.Flags().String("project-update", "", "Scope to a single project update ID")
	cmd.Flags().String("initiative-update", "", "Scope to a single initiative update ID")
}

var inboxMarkAllReadCmd = &cobra.Command{
	Use:   "mark-all-read",
	Short: "Mark all notifications as read",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.MarkAllNotificationsRead(context.Background(), time.Now()); err != nil {
			output.Error(fmt.Sprintf("Failed to mark all read: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Marked all notifications as read", plaintext, jsonOut)
	},
}

var inboxMarkAllUnreadCmd = &cobra.Command{
	Use:   "mark-all-unread",
	Short: "Mark notifications as unread",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.MarkAllNotificationsUnread(context.Background(), notificationEntityInput(cmd)); err != nil {
			output.Error(fmt.Sprintf("Failed to mark unread: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Marked notifications as unread", plaintext, jsonOut)
	},
}

var inboxSnoozeAllCmd = &cobra.Command{
	Use:   "snooze-all --until RFC3339",
	Short: "Snooze notifications",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		until, _ := cmd.Flags().GetString("until")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.SnoozeAllNotifications(context.Background(), notificationEntityInput(cmd), until); err != nil {
			output.Error(fmt.Sprintf("Failed to snooze: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Snoozed notifications until %s", until), plaintext, jsonOut)
	},
}

var inboxUnsnoozeAllCmd = &cobra.Command{
	Use:   "unsnooze-all",
	Short: "Clear notification snoozes",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UnsnoozeAllNotifications(context.Background(), notificationEntityInput(cmd), time.Now().Format(time.RFC3339)); err != nil {
			output.Error(fmt.Sprintf("Failed to unsnooze: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Cleared notification snoozes", plaintext, jsonOut)
	},
}

var inboxArchiveAllCmd = &cobra.Command{
	Use:   "archive-all",
	Short: "Archive notifications",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.ArchiveAllNotifications(context.Background(), notificationEntityInput(cmd)); err != nil {
			output.Error(fmt.Sprintf("Failed to archive notifications: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success("Archived notifications", plaintext, jsonOut)
	},
}

func init() {
	inboxCmd.AddCommand(inboxMarkAllReadCmd)
	inboxCmd.AddCommand(inboxMarkAllUnreadCmd)
	inboxCmd.AddCommand(inboxSnoozeAllCmd)
	inboxCmd.AddCommand(inboxUnsnoozeAllCmd)
	inboxCmd.AddCommand(inboxArchiveAllCmd)

	for _, c := range []*cobra.Command{inboxMarkAllUnreadCmd, inboxSnoozeAllCmd, inboxUnsnoozeAllCmd, inboxArchiveAllCmd} {
		addNotificationEntityFlags(c)
	}
	inboxSnoozeAllCmd.Flags().String("until", "", "Snooze until (RFC3339) (required)")
	_ = inboxSnoozeAllCmd.MarkFlagRequired("until")
}
