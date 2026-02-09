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

// inboxCmd represents the inbox command
var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "View Linear notifications",
	Long: `View your Linear notifications (inbox items).

Shows notifications such as issue assignments, mentions, comments,
due date reminders, and more.

Examples:
  linear-cli inbox                    # List recent notifications
  linear-cli inbox --limit 20         # List 20 notifications
  linear-cli inbox --unread           # Show only unread notifications
  linear-cli inbox --all              # Include archived notifications
  linear-cli inbox --json             # Output as JSON
  linear-cli inbox --plaintext        # Output as plaintext`,
	Run: runInboxList,
}

func runInboxList(cmd *cobra.Command, args []string) {
	plaintext := viper.GetBool("plaintext")
	jsonOut := viper.GetBool("json")

	// Get auth header
	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	// Create API client
	client := api.NewClient(authHeader)

	// Get flags
	limit, _ := cmd.Flags().GetInt("limit")
	unreadOnly, _ := cmd.Flags().GetBool("unread")
	includeArchived, _ := cmd.Flags().GetBool("all")

	// Get notifications
	notifications, err := client.GetNotifications(context.Background(), limit, "", includeArchived)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to get notifications: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	// Filter to unread if requested
	filteredNotifications := notifications.Nodes
	if unreadOnly {
		var unread []api.Notification
		for _, n := range notifications.Nodes {
			if n.ReadAt == nil {
				unread = append(unread, n)
			}
		}
		filteredNotifications = unread
	}

	// Handle empty results
	if len(filteredNotifications) == 0 {
		if jsonOut {
			output.JSON([]interface{}{})
		} else if plaintext {
			fmt.Println("No notifications found")
		} else {
			fmt.Printf("%s No notifications found\n", color.New(color.FgYellow).Sprint("📭"))
		}
		return
	}

	// Handle output
	if jsonOut {
		output.JSON(filteredNotifications)
	} else if plaintext {
		fmt.Println("Type\tIssue\tTitle\tTime\tRead")
		for _, n := range filteredNotifications {
			issueID := ""
			title := ""
			if n.Issue != nil {
				issueID = n.Issue.Identifier
				title = n.Issue.Title
			}
			readStatus := "unread"
			if n.ReadAt != nil {
				readStatus = "read"
			}
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
				formatNotificationType(n.Type),
				issueID,
				title,
				formatRelativeTime(n.CreatedAt),
				readStatus,
			)
		}
	} else {
		// Rich table output
		headers := []string{"Type", "Issue", "Title", "Time", "Status"}
		rows := [][]string{}

		for _, n := range filteredNotifications {
			notifType := formatNotificationTypeColored(n.Type)

			issueID := ""
			title := ""
			if n.Issue != nil {
				issueID = color.New(color.FgCyan).Sprint(n.Issue.Identifier)
				title = truncateString(n.Issue.Title, 50)
			}

			timeStr := color.New(color.FgWhite).Sprint(formatRelativeTime(n.CreatedAt))

			status := color.New(color.FgYellow, color.Bold).Sprint("UNREAD")
			if n.ReadAt != nil {
				status = color.New(color.FgWhite).Sprint("read")
			}

			rows = append(rows, []string{
				notifType,
				issueID,
				title,
				timeStr,
				status,
			})
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		// Summary
		unreadCount := 0
		for _, n := range filteredNotifications {
			if n.ReadAt == nil {
				unreadCount++
			}
		}
		fmt.Printf("\n%s %d notifications (%d unread)\n",
			color.New(color.FgGreen).Sprint("📬"),
			len(filteredNotifications),
			unreadCount)
	}
}

// formatNotificationType returns a human-readable notification type
func formatNotificationType(notifType string) string {
	typeMap := map[string]string{
		"issueAssignedToYou":     "assigned",
		"issueMention":           "mentioned",
		"issueNewComment":        "comment",
		"issueCommentMention":    "mentioned",
		"issueCommentReaction":   "reaction",
		"issueDue":               "due",
		"issueStatusChanged":     "status",
		"issueSubscribed":        "subscribed",
		"issuePriorityUrgent":    "urgent",
		"issueCreated":           "created",
		"issueBlocking":          "blocking",
		"issueUnblocked":         "unblocked",
		"projectUpdateCreated":   "update",
		"projectUpdateMention":   "mentioned",
	}

	if readable, ok := typeMap[notifType]; ok {
		return readable
	}
	// Convert camelCase to readable format
	return strings.ToLower(notifType)
}

// formatNotificationTypeColored returns a colored notification type
func formatNotificationTypeColored(notifType string) string {
	readable := formatNotificationType(notifType)

	switch readable {
	case "assigned":
		return color.New(color.FgMagenta).Sprint("assigned")
	case "mentioned":
		return color.New(color.FgYellow).Sprint("mentioned")
	case "comment":
		return color.New(color.FgBlue).Sprint("comment")
	case "reaction":
		return color.New(color.FgCyan).Sprint("reaction")
	case "due":
		return color.New(color.FgRed).Sprint("due")
	case "status":
		return color.New(color.FgGreen).Sprint("status")
	case "urgent":
		return color.New(color.FgRed, color.Bold).Sprint("URGENT")
	case "created":
		return color.New(color.FgGreen).Sprint("created")
	case "blocking":
		return color.New(color.FgRed).Sprint("blocking")
	case "unblocked":
		return color.New(color.FgGreen).Sprint("unblocked")
	case "update":
		return color.New(color.FgBlue).Sprint("update")
	default:
		return color.New(color.FgWhite).Sprint(readable)
	}
}

// formatRelativeTime formats a time as a relative time string
func formatRelativeTime(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2006-01-02")
	}
}

func init() {
	rootCmd.AddCommand(inboxCmd)

	// Inbox flags
	inboxCmd.Flags().IntP("limit", "l", 50, "Maximum number of notifications to return")
	inboxCmd.Flags().BoolP("unread", "u", false, "Show only unread notifications")
	inboxCmd.Flags().BoolP("all", "a", false, "Include archived notifications")
}
