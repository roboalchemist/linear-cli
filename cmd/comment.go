package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// commentCmd represents the comment command (nested under issue)
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage issue comments",
	Long: `Manage comments on Linear issues including listing and creating comments.

Examples:
  linear-cli issue comment list LIN-123                      # List comments for an issue
  linear-cli issue comment create LIN-123 --body "Fixed"     # Add a comment
  linear-cli comment list LIN-123                            # Also works (shortcut)`,
}

var commentListCmd = &cobra.Command{
	Use:     "list ISSUE-ID",
	Aliases: []string{"ls"},
	Short:   "List comments for an issue",
	Long:    `List all comments for a specific issue.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		issueID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get limit
		limit, _ := cmd.Flags().GetInt("limit")

		// Get sort option
		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				// Use empty string for Linear's default sort
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		// Get comments
		comments, err := client.GetIssueComments(context.Background(), issueID, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list comments: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(comments.Nodes)
		} else if plaintext {
			for i, comment := range comments.Nodes {
				if i > 0 {
					fmt.Println("---")
				}
				fmt.Printf("ID: %s\n", comment.ID)
				fmt.Printf("Author: %s\n", getCommentAuthor(&comment))
				fmt.Printf("Date: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
				if comment.EditedAt != nil {
					fmt.Printf("Edited: %s\n", comment.EditedAt.Format("2006-01-02 15:04:05"))
				}
				if comment.ResolvedAt != nil {
					fmt.Printf("Resolved: %s", comment.ResolvedAt.Format("2006-01-02 15:04:05"))
					if comment.ResolvingUser != nil {
						fmt.Printf(" by %s", safeUserName(comment.ResolvingUser))
					}
					fmt.Println()
				}
				if comment.ParentID != nil && *comment.ParentID != "" {
					fmt.Printf("Reply to: %s\n", *comment.ParentID)
				}
				if comment.QuotedText != nil && *comment.QuotedText != "" {
					fmt.Printf("Quoted: %s\n", *comment.QuotedText)
				}
				if comment.URL != "" {
					fmt.Printf("URL: %s\n", comment.URL)
				}
				fmt.Printf("Comment:\n%s\n", comment.Body)
			}
		} else {
			// Rich display
			if len(comments.Nodes) == 0 {
				fmt.Printf("\n%s No comments on issue %s\n",
					color.New(color.FgYellow).Sprint("ℹ"),
					color.New(color.FgCyan).Sprint(issueID))
				return
			}

			fmt.Printf("\n%s Comments on %s (%d)\n\n",
				color.New(color.FgCyan, color.Bold).Sprint("💬"),
				color.New(color.FgCyan).Sprint(issueID),
				len(comments.Nodes))

			for i, comment := range comments.Nodes {
				if i > 0 {
					fmt.Println(strings.Repeat("─", 50))
				}

				// Header with author and time
				timeAgo := formatTimeAgo(comment.CreatedAt)
				authorDisplay := getCommentAuthor(&comment)

				// Show thread indicator for replies
				if comment.ParentID != nil && *comment.ParentID != "" {
					fmt.Printf("%s ", color.New(color.FgWhite, color.Faint).Sprint("↳"))
				}

				fmt.Printf("%s %s %s",
					color.New(color.FgCyan, color.Bold).Sprint(authorDisplay),
					color.New(color.FgWhite, color.Faint).Sprint("•"),
					color.New(color.FgWhite, color.Faint).Sprint(timeAgo))

				// Show edited indicator
				if comment.EditedAt != nil {
					fmt.Printf(" %s", color.New(color.FgWhite, color.Faint).Sprint("(edited)"))
				}

				// Show resolved indicator
				if comment.ResolvedAt != nil {
					fmt.Printf(" %s", color.New(color.FgGreen).Sprint("✓ resolved"))
				}

				fmt.Println()

				// Show quoted text if present
				if comment.QuotedText != nil && *comment.QuotedText != "" {
					fmt.Printf("%s %s\n",
						color.New(color.FgWhite, color.Faint).Sprint("│"),
						color.New(color.FgWhite, color.Faint).Sprint(*comment.QuotedText))
				}

				// Comment body
				fmt.Printf("\n%s\n\n", comment.Body)
			}
		}
	},
}

var commentCreateCmd = &cobra.Command{
	Use:     "create ISSUE-ID",
	Aliases: []string{"add", "new"},
	Short:   "Create a comment on an issue",
	Long: `Add a new comment to a specific issue.

The comment body can be provided inline via --body or read from a markdown file via --body-file.
Use --body-file - to read from stdin.

Threaded comments:
  Use --parent to reply to an existing comment, creating a threaded conversation.

Examples:
  linear-cli issue comment create LIN-123 --body "Fixed the bug"
  linear-cli issue comment create LIN-123 --body-file comment.md
  linear-cli issue comment create LIN-123 --body "Reply" --parent COMMENT-UUID
  linear-cli issue comment create LIN-123 --body "Note" --quoted-text "Original text"
  cat notes.md | linear-cli issue comment create LIN-123 --body-file -`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		issueID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Resolve comment body from --body or --body-file
		bodyFlag, _ := cmd.Flags().GetString("body")
		filePath, _ := cmd.Flags().GetString("body-file")
		body, err := resolveBodyFromFlags(bodyFlag, cmd.Flags().Changed("body"), filePath, "body", "body-file")
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}
		if body == "" {
			output.Error("Comment body is required (--body or --body-file)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Build options
		opts := &api.CommentCreateOptions{}
		if parentID, _ := cmd.Flags().GetString("parent"); parentID != "" {
			opts.ParentID = parentID
		}
		if quotedText, _ := cmd.Flags().GetString("quoted-text"); quotedText != "" {
			opts.QuotedText = quotedText
		}
		if doNotSubscribe, _ := cmd.Flags().GetBool("do-not-subscribe"); doNotSubscribe {
			opts.DoNotSubscribe = true
		}

		// Create comment
		comment, err := client.CreateComment(context.Background(), issueID, body, opts)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(comment)
		} else if plaintext {
			fmt.Printf("Created comment on %s\n", issueID)
			fmt.Printf("ID: %s\n", comment.ID)
			fmt.Printf("Author: %s\n", getCommentAuthor(comment))
			fmt.Printf("Date: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
			if comment.URL != "" {
				fmt.Printf("URL: %s\n", comment.URL)
			}
			if comment.ParentID != nil && *comment.ParentID != "" {
				fmt.Printf("Reply to: %s\n", *comment.ParentID)
			}
		} else {
			fmt.Printf("%s Added comment to %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(issueID))
			if comment.ParentID != nil && *comment.ParentID != "" {
				fmt.Printf("   %s Reply to comment %s\n",
					color.New(color.FgWhite, color.Faint).Sprint("↳"),
					color.New(color.FgWhite, color.Faint).Sprint(*comment.ParentID))
			}
			fmt.Printf("\n%s\n", comment.Body)
		}
	},
}

// safeUserName returns the user's name, falling back to email, then "System"
func safeUserName(user *api.User) string {
	if user == nil {
		return "System"
	}
	if name := strings.TrimSpace(user.Name); name != "" {
		return name
	}
	if email := strings.TrimSpace(user.Email); email != "" {
		return email
	}
	return "System"
}

// getCommentAuthor returns the author of a comment, checking user, externalUser, and botActor
func getCommentAuthor(comment *api.Comment) string {
	if comment.User != nil {
		return safeUserName(comment.User)
	}
	if comment.ExternalUser != nil {
		if name := strings.TrimSpace(comment.ExternalUser.Name); name != "" {
			return name + " (external)"
		}
		if email := strings.TrimSpace(comment.ExternalUser.Email); email != "" {
			return email + " (external)"
		}
	}
	if comment.BotActor != nil {
		if comment.BotActor.UserDisplayName != nil && *comment.BotActor.UserDisplayName != "" {
			return *comment.BotActor.UserDisplayName + " (bot)"
		}
		if comment.BotActor.Name != nil && *comment.BotActor.Name != "" {
			return *comment.BotActor.Name + " (bot)"
		}
		return "Bot"
	}
	return "System"
}

// formatTimeAgo formats a time as a human-readable "time ago" string
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(duration.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

var commentUpdateCmd = &cobra.Command{
	Use:     "update COMMENT-ID",
	Aliases: []string{"edit"},
	Short:   "Update a comment",
	Long: `Update an existing comment's body, resolution status, or other fields.

The comment body can be provided inline via --body or read from a markdown file via --body-file.
Use --body-file - to read from stdin.

Resolution:
  Use --resolve to mark a comment as resolved (e.g., feedback addressed).
  Use --unresolve to clear the resolution status.

Examples:
  linear-cli issue comment update COMMENT-ID --body "Updated text"
  linear-cli issue comment update COMMENT-ID --body-file updated.md
  linear-cli issue comment update COMMENT-ID --resolve
  linear-cli issue comment update COMMENT-ID --unresolve
  linear-cli issue comment update COMMENT-ID --quoted-text "Referenced text"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		commentID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Build options
		opts := &api.CommentUpdateOptions{}
		hasChanges := false

		// Resolve comment body from --body or --body-file
		bodyFlag, _ := cmd.Flags().GetString("body")
		filePath, _ := cmd.Flags().GetString("body-file")
		if cmd.Flags().Changed("body") || cmd.Flags().Changed("body-file") {
			body, err := resolveBodyFromFlags(bodyFlag, cmd.Flags().Changed("body"), filePath, "body", "body-file")
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			opts.Body = &body
			hasChanges = true
		}

		if quotedText, _ := cmd.Flags().GetString("quoted-text"); cmd.Flags().Changed("quoted-text") {
			opts.QuotedText = &quotedText
			hasChanges = true
		}

		resolve, _ := cmd.Flags().GetBool("resolve")
		unresolve, _ := cmd.Flags().GetBool("unresolve")

		if resolve && unresolve {
			output.Error("Cannot use both --resolve and --unresolve", plaintext, jsonOut)
			os.Exit(1)
		}

		if resolve {
			// Get current user ID to set as resolver
			viewer, err := client.GetViewer(context.Background())
			if err != nil {
				output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			opts.ResolvingUserID = &viewer.ID
			hasChanges = true
		}

		if unresolve {
			emptyStr := ""
			opts.ResolvingUserID = &emptyStr
			hasChanges = true
		}

		if doNotSubscribe, _ := cmd.Flags().GetBool("do-not-subscribe"); doNotSubscribe {
			opts.DoNotSubscribe = true
		}

		if !hasChanges {
			output.Error("No changes specified. Use --body, --resolve, --unresolve, or --quoted-text", plaintext, jsonOut)
			os.Exit(1)
		}

		comment, err := client.UpdateComment(context.Background(), commentID, opts)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(comment)
		} else if plaintext {
			fmt.Println("Updated comment")
			fmt.Printf("ID: %s\n", comment.ID)
			if comment.EditedAt != nil {
				fmt.Printf("Edited at: %s\n", comment.EditedAt.Format("2006-01-02 15:04:05"))
			}
			if comment.ResolvedAt != nil {
				fmt.Printf("Resolved at: %s\n", comment.ResolvedAt.Format("2006-01-02 15:04:05"))
				if comment.ResolvingUser != nil {
					fmt.Printf("Resolved by: %s\n", safeUserName(comment.ResolvingUser))
				}
			}
		} else {
			output.Success("Updated comment", plaintext, jsonOut)
			if comment.ResolvedAt != nil {
				fmt.Printf("   %s Resolved by %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan).Sprint(safeUserName(comment.ResolvingUser)))
			}
		}
	},
}

var commentDeleteCmd = &cobra.Command{
	Use:     "delete COMMENT-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a comment",
	Long:    `Delete an existing comment.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		commentID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		err = client.DeleteComment(context.Background(), commentID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted comment", plaintext, jsonOut)
	},
}

func init() {
	// Primary home: nested under issue command
	issueCmd.AddCommand(commentCmd)
	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentCreateCmd)
	commentCmd.AddCommand(commentUpdateCmd)
	commentCmd.AddCommand(commentDeleteCmd)

	// List command flags
	commentListCmd.Flags().IntP("limit", "l", 50, "Maximum number of comments to return")
	commentListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")

	// Create command flags
	commentCreateCmd.Flags().StringP("body", "b", "", "Comment body (required unless --body-file is used)")
	commentCreateCmd.Flags().String("body-file", "", "Read body from a markdown file (use - for stdin)")
	commentCreateCmd.Flags().String("parent", "", "Parent comment ID (for threaded replies)")
	commentCreateCmd.Flags().String("quoted-text", "", "Text being quoted or referenced")
	commentCreateCmd.Flags().Bool("do-not-subscribe", false, "Don't subscribe to the issue after commenting")

	// Update command flags
	commentUpdateCmd.Flags().StringP("body", "b", "", "New comment body")
	commentUpdateCmd.Flags().String("body-file", "", "Read body from a markdown file (use - for stdin)")
	commentUpdateCmd.Flags().String("quoted-text", "", "Text being quoted or referenced")
	commentUpdateCmd.Flags().Bool("resolve", false, "Mark the comment as resolved")
	commentUpdateCmd.Flags().Bool("unresolve", false, "Clear the resolution status")
	commentUpdateCmd.Flags().Bool("do-not-subscribe", false, "Don't subscribe to the issue")
}
