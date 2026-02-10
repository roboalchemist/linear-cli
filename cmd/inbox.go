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
		fmt.Println("Type\tIssue\tTitle\tActor\tTeam\tTime\tRead")
		for _, n := range filteredNotifications {
			issueID := ""
			title := n.Title
			if n.Issue != nil {
				issueID = n.Issue.Identifier
				if title == "" {
					title = n.Issue.Title
				}
			} else if n.Project != nil {
				issueID = n.Project.Name
			}
			actorName := ""
			if n.Actor != nil {
				actorName = n.Actor.Name
			}
			teamKey := ""
			if n.Team != nil {
				teamKey = n.Team.Key
			}
			readStatus := "unread"
			if n.ReadAt != nil {
				readStatus = "read"
			}
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				formatNotificationType(n.Type),
				issueID,
				title,
				actorName,
				teamKey,
				formatRelativeTime(n.CreatedAt),
				readStatus,
			)
		}
	} else {
		// Rich table output
		headers := []string{"Type", "Issue", "Title", "Actor", "Team", "Time", "Status"}
		rows := [][]string{}

		for _, n := range filteredNotifications {
			notifType := formatNotificationTypeColored(n.Type)

			issueID := ""
			title := truncateString(n.Title, 40)
			if n.Issue != nil {
				issueID = color.New(color.FgCyan).Sprint(n.Issue.Identifier)
				if title == "" {
					title = truncateString(n.Issue.Title, 40)
				}
			} else if n.Project != nil {
				issueID = color.New(color.FgMagenta).Sprint(n.Project.Name)
			}

			actorName := ""
			if n.Actor != nil {
				actorName = color.New(color.FgWhite).Sprint(n.Actor.Name)
			}

			teamKey := ""
			if n.Team != nil {
				teamKey = color.New(color.FgBlue).Sprint(n.Team.Key)
			}

			timeStr := color.New(color.FgWhite).Sprint(formatRelativeTime(n.CreatedAt))

			status := color.New(color.FgYellow, color.Bold).Sprint("UNREAD")
			if n.ReadAt != nil {
				status = color.New(color.FgWhite).Sprint("read")
			}
			if n.SnoozedUntilAt != nil {
				status = color.New(color.FgCyan).Sprint("snoozed")
			}

			rows = append(rows, []string{
				notifType,
				issueID,
				title,
				actorName,
				teamKey,
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
		snoozedCount := 0
		for _, n := range filteredNotifications {
			if n.ReadAt == nil {
				unreadCount++
			}
			if n.SnoozedUntilAt != nil {
				snoozedCount++
			}
		}
		summaryParts := []string{fmt.Sprintf("%d unread", unreadCount)}
		if snoozedCount > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d snoozed", snoozedCount))
		}
		fmt.Printf("\n%s %d notifications (%s)\n",
			color.New(color.FgGreen).Sprint("📬"),
			len(filteredNotifications),
			strings.Join(summaryParts, ", "))
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

// inboxReadCmd marks a notification as read
var inboxReadCmd = &cobra.Command{
	Use:   "read <notification-id>",
	Short: "Mark a notification as read",
	Long: `Mark a notification as read.

Use --all to mark all notifications as read.

Examples:
  linear-cli inbox read abc123           # Mark specific notification as read
  linear-cli inbox read --all            # Mark all notifications as read`,
	Run: runInboxRead,
}

func runInboxRead(cmd *cobra.Command, args []string) {
	plaintext := viper.GetBool("plaintext")
	jsonOut := viper.GetBool("json")

	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	client := api.NewClient(authHeader)

	markAll, _ := cmd.Flags().GetBool("all")

	if markAll {
		// Mark all as read
		err = client.MarkAllNotificationsRead(context.Background(), time.Now())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to mark notifications as read: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "message": "All notifications marked as read"})
		} else if plaintext {
			fmt.Println("All notifications marked as read")
		} else {
			fmt.Printf("%s All notifications marked as read\n", color.New(color.FgGreen).Sprint("✓"))
		}
		return
	}

	if len(args) == 0 {
		output.Error("Notification ID required (or use --all)", plaintext, jsonOut)
		os.Exit(1)
	}

	notificationID := args[0]
	now := time.Now()
	input := api.NotificationUpdateInput{
		ReadAt: &now,
	}

	notification, err := client.UpdateNotification(context.Background(), notificationID, input)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to mark notification as read: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	if jsonOut {
		output.JSON(notification)
	} else if plaintext {
		fmt.Printf("Notification %s marked as read\n", notificationID)
	} else {
		fmt.Printf("%s Notification %s marked as read\n",
			color.New(color.FgGreen).Sprint("✓"),
			color.New(color.FgCyan).Sprint(notificationID))
	}
}

// inboxUnreadCmd marks a notification as unread
var inboxUnreadCmd = &cobra.Command{
	Use:   "unread <notification-id>",
	Short: "Mark a notification as unread",
	Long: `Mark a notification as unread.

Examples:
  linear-cli inbox unread abc123         # Mark notification as unread`,
	Args: cobra.ExactArgs(1),
	Run:  runInboxUnread,
}

func runInboxUnread(cmd *cobra.Command, args []string) {
	plaintext := viper.GetBool("plaintext")
	jsonOut := viper.GetBool("json")

	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	client := api.NewClient(authHeader)
	notificationID := args[0]

	input := api.NotificationUpdateInput{
		ReadAt: nil, // Setting to nil marks as unread
	}

	notification, err := client.UpdateNotification(context.Background(), notificationID, input)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to mark notification as unread: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	if jsonOut {
		output.JSON(notification)
	} else if plaintext {
		fmt.Printf("Notification %s marked as unread\n", notificationID)
	} else {
		fmt.Printf("%s Notification %s marked as unread\n",
			color.New(color.FgGreen).Sprint("✓"),
			color.New(color.FgCyan).Sprint(notificationID))
	}
}

// inboxSnoozeCmd snoozes a notification
var inboxSnoozeCmd = &cobra.Command{
	Use:   "snooze <notification-id> <duration>",
	Short: "Snooze a notification",
	Long: `Snooze a notification for a specified duration.

Duration can be specified as:
  - 1h, 2h, etc. for hours
  - 1d, 2d, etc. for days
  - 1w for a week
  - tomorrow (snooze until 9am tomorrow)

Examples:
  linear-cli inbox snooze abc123 1h       # Snooze for 1 hour
  linear-cli inbox snooze abc123 1d       # Snooze for 1 day
  linear-cli inbox snooze abc123 tomorrow # Snooze until tomorrow`,
	Args: cobra.ExactArgs(2),
	Run:  runInboxSnooze,
}

func runInboxSnooze(cmd *cobra.Command, args []string) {
	plaintext := viper.GetBool("plaintext")
	jsonOut := viper.GetBool("json")

	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	client := api.NewClient(authHeader)
	notificationID := args[0]
	durationStr := args[1]

	// Parse the duration
	var snoozeUntil time.Time
	switch {
	case durationStr == "tomorrow":
		// Snooze until 9am tomorrow
		now := time.Now()
		tomorrow := now.AddDate(0, 0, 1)
		snoozeUntil = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 9, 0, 0, 0, now.Location())
	case strings.HasSuffix(durationStr, "w"):
		weeks, err := parseIntFromSuffix(durationStr, "w")
		if err != nil {
			output.Error(fmt.Sprintf("Invalid duration: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		snoozeUntil = time.Now().AddDate(0, 0, weeks*7)
	case strings.HasSuffix(durationStr, "d"):
		days, err := parseIntFromSuffix(durationStr, "d")
		if err != nil {
			output.Error(fmt.Sprintf("Invalid duration: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		snoozeUntil = time.Now().AddDate(0, 0, days)
	case strings.HasSuffix(durationStr, "h"):
		hours, err := parseIntFromSuffix(durationStr, "h")
		if err != nil {
			output.Error(fmt.Sprintf("Invalid duration: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		snoozeUntil = time.Now().Add(time.Duration(hours) * time.Hour)
	default:
		// Try parsing as a Go duration
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			output.Error(fmt.Sprintf("Invalid duration format: %s", durationStr), plaintext, jsonOut)
			os.Exit(1)
		}
		snoozeUntil = time.Now().Add(duration)
	}

	input := api.NotificationUpdateInput{
		SnoozedUntilAt: &snoozeUntil,
	}

	notification, err := client.UpdateNotification(context.Background(), notificationID, input)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to snooze notification: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	if jsonOut {
		output.JSON(notification)
	} else if plaintext {
		fmt.Printf("Notification %s snoozed until %s\n", notificationID, snoozeUntil.Format(time.RFC3339))
	} else {
		fmt.Printf("%s Notification %s snoozed until %s\n",
			color.New(color.FgGreen).Sprint("✓"),
			color.New(color.FgCyan).Sprint(notificationID),
			color.New(color.FgYellow).Sprint(snoozeUntil.Format("Jan 2 at 3:04 PM")))
	}
}

// parseIntFromSuffix parses an integer from a string with a suffix (e.g., "2d" -> 2)
func parseIntFromSuffix(s, suffix string) (int, error) {
	numStr := strings.TrimSuffix(s, suffix)
	var num int
	_, err := fmt.Sscanf(numStr, "%d", &num)
	if err != nil {
		return 0, fmt.Errorf("could not parse number from %s", s)
	}
	return num, nil
}

// inboxArchiveCmd archives a notification
var inboxArchiveCmd = &cobra.Command{
	Use:   "archive <notification-id>",
	Short: "Archive a notification",
	Long: `Archive a notification to remove it from your inbox.

Examples:
  linear-cli inbox archive abc123        # Archive notification`,
	Args: cobra.ExactArgs(1),
	Run:  runInboxArchive,
}

func runInboxArchive(cmd *cobra.Command, args []string) {
	plaintext := viper.GetBool("plaintext")
	jsonOut := viper.GetBool("json")

	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	client := api.NewClient(authHeader)
	notificationID := args[0]

	err = client.ArchiveNotification(context.Background(), notificationID)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to archive notification: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	if jsonOut {
		output.JSON(map[string]interface{}{"success": true, "id": notificationID})
	} else if plaintext {
		fmt.Printf("Notification %s archived\n", notificationID)
	} else {
		fmt.Printf("%s Notification %s archived\n",
			color.New(color.FgGreen).Sprint("✓"),
			color.New(color.FgCyan).Sprint(notificationID))
	}
}

// inboxUnarchiveCmd unarchives a notification
var inboxUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <notification-id>",
	Short: "Unarchive a notification",
	Long: `Unarchive a notification to return it to your inbox.

Examples:
  linear-cli inbox unarchive abc123      # Unarchive notification`,
	Args: cobra.ExactArgs(1),
	Run:  runInboxUnarchive,
}

func runInboxUnarchive(cmd *cobra.Command, args []string) {
	plaintext := viper.GetBool("plaintext")
	jsonOut := viper.GetBool("json")

	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	client := api.NewClient(authHeader)
	notificationID := args[0]

	err = client.UnarchiveNotification(context.Background(), notificationID)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to unarchive notification: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}

	if jsonOut {
		output.JSON(map[string]interface{}{"success": true, "id": notificationID})
	} else if plaintext {
		fmt.Printf("Notification %s unarchived\n", notificationID)
	} else {
		fmt.Printf("%s Notification %s unarchived\n",
			color.New(color.FgGreen).Sprint("✓"),
			color.New(color.FgCyan).Sprint(notificationID))
	}
}

func init() {
	rootCmd.AddCommand(inboxCmd)

	// Inbox list flags
	inboxCmd.Flags().IntP("limit", "l", 50, "Maximum number of notifications to return")
	inboxCmd.Flags().BoolP("unread", "u", false, "Show only unread notifications")
	inboxCmd.Flags().BoolP("all", "a", false, "Include archived notifications")

	// Subcommands
	inboxCmd.AddCommand(inboxReadCmd)
	inboxReadCmd.Flags().BoolP("all", "a", false, "Mark all notifications as read")

	inboxCmd.AddCommand(inboxUnreadCmd)
	inboxCmd.AddCommand(inboxSnoozeCmd)
	inboxCmd.AddCommand(inboxArchiveCmd)
	inboxCmd.AddCommand(inboxUnarchiveCmd)
}
