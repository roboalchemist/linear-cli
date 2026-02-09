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
				fmt.Printf("Author: %s\n", safeUserName(comment.User))
				fmt.Printf("Date: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("Comment:\n%s\n", comment.Body)
			}
		} else {
			// Rich display
			if len(comments.Nodes) == 0 {
				fmt.Printf("\n%s No comments on issue %s\n",
					color.New(color.FgYellow).Sprint("ℹ️"),
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
				fmt.Printf("%s %s %s\n",
					color.New(color.FgCyan, color.Bold).Sprint(safeUserName(comment.User)),
					color.New(color.FgWhite, color.Faint).Sprint("•"),
					color.New(color.FgWhite, color.Faint).Sprint(timeAgo))

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

The comment body can be provided inline via --body or read from a markdown file via --file.
Use --file - to read from stdin.

Examples:
  linear-cli issue comment create LIN-123 --body "Fixed the bug"
  linear-cli issue comment create LIN-123 --file comment.md
  cat notes.md | linear-cli issue comment create LIN-123 --file -`,
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

		// Resolve comment body from --body or --file
		bodyFlag, _ := cmd.Flags().GetString("body")
		filePath, _ := cmd.Flags().GetString("file")
		body, err := resolveBodyFromFlags(bodyFlag, cmd.Flags().Changed("body"), filePath, "body")
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}
		if body == "" {
			output.Error("Comment body is required (--body or --file)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Create comment
		comment, err := client.CreateComment(context.Background(), issueID, body)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(comment)
		} else if plaintext {
			fmt.Printf("Created comment on %s\n", issueID)
			fmt.Printf("Author: %s\n", safeUserName(comment.User))
			fmt.Printf("Date: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%s Added comment to %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(issueID))
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
	Long: `Update the body of an existing comment.

The comment body can be provided inline via --body or read from a markdown file via --file.
Use --file - to read from stdin.

Examples:
  linear-cli issue comment update COMMENT-ID --body "Updated text"
  linear-cli issue comment update COMMENT-ID --file updated.md`,
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

		// Resolve comment body from --body or --file
		bodyFlag, _ := cmd.Flags().GetString("body")
		filePath, _ := cmd.Flags().GetString("file")
		body, err := resolveBodyFromFlags(bodyFlag, cmd.Flags().Changed("body"), filePath, "body")
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}
		if body == "" {
			output.Error("Comment body is required (--body or --file)", plaintext, jsonOut)
			os.Exit(1)
		}

		comment, err := client.UpdateComment(context.Background(), commentID, body)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(comment)
		} else {
			output.Success("Updated comment", plaintext, jsonOut)
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
	commentCreateCmd.Flags().StringP("body", "b", "", "Comment body (required unless --file is used)")
	commentCreateCmd.Flags().StringP("file", "f", "", "Read body from a markdown file (use - for stdin)")

	// Update command flags
	commentUpdateCmd.Flags().StringP("body", "b", "", "New comment body (required unless --file is used)")
	commentUpdateCmd.Flags().StringP("file", "f", "", "Read body from a markdown file (use - for stdin)")
}
