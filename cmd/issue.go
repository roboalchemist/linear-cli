package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/pkg/auth"
	"github.com/dorkitude/linctl/pkg/output"
	"github.com/dorkitude/linctl/pkg/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Linear issues",
	Long: `Create, list, update, and manage Linear issues.

Examples:
  linctl issue list --assignee me --state "In Progress"
  linctl issue ls -a me -s "In Progress"
  linctl issue list --include-completed  # Show all issues including completed
  linctl issue list --newer-than 3_weeks_ago  # Show issues from last 3 weeks
  linctl issue get LIN-123
  linctl issue create --title "Bug fix" --team ENG`,
}

var issueListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issues",
	Long:    `List Linear issues with optional filtering.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Build filter from flags
		filter := buildIssueFilter(cmd)

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 50
		}

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

		issues, err := client.GetIssues(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(issues.Nodes) == 0 {
			output.Info("No issues found", plaintext, jsonOut)
			return
		}

		// For JSON output, show raw data
		if jsonOut {
			output.JSON(issues.Nodes)
			return
		}

		// Prepare table data
		headers := []string{"Identifier", "Title", "State", "Assignee", "Team", "Priority", "Created", "Updated"}
		rows := make([][]string, len(issues.Nodes))

		for i, issue := range issues.Nodes {
			assignee := "Unassigned"
			if issue.Assignee != nil {
				assignee = issue.Assignee.Name
			}

			team := ""
			if issue.Team != nil {
				team = issue.Team.Key
			}

			state := ""
			if issue.State != nil {
				state = issue.State.Name
			}

			priority := priorityToString(issue.Priority)

			// Apply colors if not in plaintext/json mode
			if !plaintext && !jsonOut {
				// Color state based on type
				if issue.State != nil {
					var stateColor *color.Color
					switch issue.State.Type {
					case "triage":
						stateColor = color.New(color.FgMagenta)
					case "backlog":
						stateColor = color.New(color.FgCyan)
					case "unstarted":
						stateColor = color.New(color.FgWhite)
					case "started":
						stateColor = color.New(color.FgBlue)
					case "completed":
						stateColor = color.New(color.FgGreen)
					case "canceled":
						stateColor = color.New(color.FgRed)
					default:
						stateColor = color.New(color.FgWhite)
					}
					state = stateColor.Sprint(state)
				}

				// Color priority
				var priorityColor *color.Color
				switch issue.Priority {
				case 0:
					priorityColor = color.New(color.FgWhite, color.Faint)
				case 1:
					priorityColor = color.New(color.FgRed, color.Bold)
				case 2:
					priorityColor = color.New(color.FgRed)
				case 3:
					priorityColor = color.New(color.FgYellow)
				case 4:
					priorityColor = color.New(color.FgBlue)
				default:
					priorityColor = color.New(color.FgWhite)
				}
				priority = priorityColor.Sprint(priority)

				// Color unassigned in yellow
				if issue.Assignee == nil {
					assignee = color.New(color.FgYellow).Sprint(assignee)
				}
			}

			rows[i] = []string{
				issue.Identifier,
				truncateString(issue.Title, 40),
				state,
				assignee,
				team,
				priority,
				issue.CreatedAt.Format("2006-01-02"),
				issue.UpdatedAt.Format("2006-01-02"),
			}
		}

		tableData := output.TableData{
			Headers: headers,
			Rows:    rows,
		}

		output.Table(tableData, plaintext, jsonOut)

		// Show summary count like project list does
		if !plaintext && !jsonOut {
			fmt.Printf("\n%s %d issues\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(issues.Nodes))
		}

		if !plaintext && !jsonOut && issues.PageInfo.HasNextPage {
			fmt.Printf("%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

var issueGetCmd = &cobra.Command{
	Use:     "get [issue-id]",
	Aliases: []string{"show"},
	Short:   "Get issue details",
	Long:    `Get detailed information about a specific issue.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		issue, err := client.GetIssue(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
			return
		}

		if plaintext {
			fmt.Printf("ID: %s\n", issue.Identifier)
			fmt.Printf("Title: %s\n", issue.Title)
			if issue.Description != "" {
				fmt.Printf("Description: %s\n", issue.Description)
			}
			if issue.State != nil {
				fmt.Printf("State: %s\n", issue.State.Name)
			}
			if issue.Assignee != nil {
				fmt.Printf("Assignee: %s\n", issue.Assignee.Name)
			}
			if issue.Team != nil {
				fmt.Printf("Team: %s\n", issue.Team.Name)
			}
			fmt.Printf("Priority: %s\n", priorityToString(issue.Priority))
			if issue.Project != nil {
				fmt.Printf("Project: %s\n", issue.Project.Name)
			}
			if issue.Cycle != nil {
				fmt.Printf("Cycle: %s\n", issue.Cycle.Name)
			}
			if issue.DueDate != nil && *issue.DueDate != "" {
				fmt.Printf("Due Date: %s\n", *issue.DueDate)
			}
			if issue.BranchName != "" {
				fmt.Printf("Git Branch: %s\n", issue.BranchName)
			}
			fmt.Printf("Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("URL: %s\n", issue.URL)
			return
		}

		// Rich display
		fmt.Printf("%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
			color.New(color.FgWhite, color.Bold).Sprint(issue.Title))

		if issue.Description != "" {
			fmt.Printf("\n%s\n", issue.Description)
		}

		fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Details:"))

		if issue.State != nil {
			stateStr := issue.State.Name
			if issue.State.Type == "completed" && issue.CompletedAt != nil {
				stateStr += fmt.Sprintf(" (%s)", issue.CompletedAt.Format("2006-01-02"))
			}
			fmt.Printf("State: %s\n",
				color.New(color.FgGreen).Sprint(stateStr))
		}

		if issue.Assignee != nil {
			fmt.Printf("Assignee: %s\n",
				color.New(color.FgCyan).Sprint(issue.Assignee.Name))
		} else {
			fmt.Printf("Assignee: %s\n",
				color.New(color.FgRed).Sprint("Unassigned"))
		}

		if issue.Team != nil {
			fmt.Printf("Team: %s\n",
				color.New(color.FgMagenta).Sprint(issue.Team.Name))
		}

		fmt.Printf("Priority: %s\n", priorityToString(issue.Priority))

		// Show project and cycle info
		if issue.Project != nil {
			fmt.Printf("Project: %s (%s)\n",
				color.New(color.FgBlue).Sprint(issue.Project.Name),
				color.New(color.FgWhite, color.Faint).Sprintf("%.0f%%", issue.Project.Progress*100))
		}

		if issue.Cycle != nil {
			fmt.Printf("Cycle: %s\n",
				color.New(color.FgMagenta).Sprint(issue.Cycle.Name))
		}

		fmt.Printf("Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))

		if issue.DueDate != nil && *issue.DueDate != "" {
			fmt.Printf("Due Date: %s\n",
				color.New(color.FgYellow).Sprint(*issue.DueDate))
		}

		if issue.SnoozedUntilAt != nil {
			fmt.Printf("Snoozed Until: %s\n",
				color.New(color.FgYellow).Sprint(issue.SnoozedUntilAt.Format("2006-01-02 15:04:05")))
		}

		// Show git branch if available
		if issue.BranchName != "" {
			fmt.Printf("Git Branch: %s\n",
				color.New(color.FgGreen).Sprint(issue.BranchName))
		}

		// Show URL
		if issue.URL != "" {
			fmt.Printf("URL: %s\n",
				color.New(color.FgBlue, color.Underline).Sprint(issue.URL))
		}

		// Show parent issue if this is a sub-issue
		if issue.Parent != nil {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Parent Issue:"))
			fmt.Printf("  %s %s\n",
				color.New(color.FgCyan).Sprint(issue.Parent.Identifier),
				issue.Parent.Title)
		}

		// Show sub-issues if any
		if issue.Children != nil && len(issue.Children.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Sub-issues:"))
			for _, child := range issue.Children.Nodes {
				stateIcon := "○"
				if child.State != nil {
					switch child.State.Type {
					case "completed", "done":
						stateIcon = color.New(color.FgGreen).Sprint("✓")
					case "started", "in_progress":
						stateIcon = color.New(color.FgBlue).Sprint("◐")
					case "canceled":
						stateIcon = color.New(color.FgRed).Sprint("✗")
					}
				}

				assignee := "Unassigned"
				if child.Assignee != nil {
					assignee = child.Assignee.Name
				}

				fmt.Printf("  %s %s %s (%s)\n",
					stateIcon,
					color.New(color.FgCyan).Sprint(child.Identifier),
					child.Title,
					color.New(color.FgWhite, color.Faint).Sprint(assignee))
			}
		}

		// Show attachments if any
		if issue.Attachments != nil && len(issue.Attachments.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Attachments:"))
			for _, attachment := range issue.Attachments.Nodes {
				fmt.Printf("  📎 %s - %s\n",
					attachment.Title,
					color.New(color.FgBlue, color.Underline).Sprint(attachment.URL))
			}
		}

		// Show recent comments if any
		if issue.Comments != nil && len(issue.Comments.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Recent Comments:"))
			for _, comment := range issue.Comments.Nodes {
				fmt.Printf("  💬 %s - %s\n",
					color.New(color.FgCyan).Sprint(comment.User.Name),
					color.New(color.FgWhite, color.Faint).Sprint(comment.CreatedAt.Format("2006-01-02 15:04")))
				// Show first line of comment
				lines := strings.Split(comment.Body, "\n")
				if len(lines) > 0 && lines[0] != "" {
					preview := lines[0]
					if len(preview) > 60 {
						preview = preview[:57] + "..."
					}
					fmt.Printf("     %s\n", preview)
				}
			}
			fmt.Printf("\n  %s Use 'linctl comment list %s' to see all comments\n",
				color.New(color.FgWhite, color.Faint).Sprint("→"),
				issue.Identifier)
		}
	},
}

func buildIssueFilter(cmd *cobra.Command) map[string]interface{} {
	filter := make(map[string]interface{})

	if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
		if assignee == "me" {
			// We'll need to get the current user's ID
			// For now, we'll use a special marker
			filter["assignee"] = map[string]interface{}{"isMe": map[string]interface{}{"eq": true}}
		} else {
			filter["assignee"] = map[string]interface{}{"email": map[string]interface{}{"eq": assignee}}
		}
	}

	state, _ := cmd.Flags().GetString("state")
	if state != "" {
		filter["state"] = map[string]interface{}{"name": map[string]interface{}{"eq": state}}
	} else {
		// Only filter out completed issues if no specific state is requested
		includeCompleted, _ := cmd.Flags().GetBool("include-completed")
		if !includeCompleted {
			// Filter out completed and canceled states
			filter["state"] = map[string]interface{}{
				"type": map[string]interface{}{
					"nin": []string{"completed", "canceled"},
				},
			}
		}
	}

	if team, _ := cmd.Flags().GetString("team"); team != "" {
		filter["team"] = map[string]interface{}{"key": map[string]interface{}{"eq": team}}
	}

	if priority, _ := cmd.Flags().GetInt("priority"); priority != -1 {
		filter["priority"] = map[string]interface{}{"eq": priority}
	}

	// Handle newer-than filter
	newerThan, _ := cmd.Flags().GetString("newer-than")
	createdAt, err := utils.ParseTimeExpression(newerThan)
	if err != nil {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		output.Error(fmt.Sprintf("Invalid newer-than value: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}
	if createdAt != "" {
		filter["createdAt"] = map[string]interface{}{"gte": createdAt}
	}

	return filter
}

func priorityToString(priority int) string {
	switch priority {
	case 0:
		return "None"
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Normal"
	case 4:
		return "Low"
	default:
		return "Unknown"
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

var issueAssignCmd = &cobra.Command{
	Use:   "assign [issue-id]",
	Short: "Assign issue to yourself",
	Long:  `Assign an issue to yourself.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Get current user
		viewer, err := client.GetViewer(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Update issue with assignee
		input := map[string]interface{}{
			"assigneeId": viewer.ID,
		}

		issue, err := client.UpdateIssue(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to assign issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
		} else if plaintext {
			fmt.Printf("Assigned %s to %s\n", issue.Identifier, viewer.Name)
		} else {
			fmt.Printf("%s Assigned %s to %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
				color.New(color.FgCyan).Sprint(viewer.Name))
		}
	},
}

var issueCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new issue",
	Long:    `Create a new issue in Linear.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Get flags
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		teamKey, _ := cmd.Flags().GetString("team")
		priority, _ := cmd.Flags().GetInt("priority")
		assignToMe, _ := cmd.Flags().GetBool("assign-me")

		if title == "" {
			output.Error("Title is required (--title)", plaintext, jsonOut)
			os.Exit(1)
		}

		if teamKey == "" {
			output.Error("Team is required (--team)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Get team ID from key
		team, err := client.GetTeam(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Build input
		input := map[string]interface{}{
			"title":  title,
			"teamId": team.ID,
		}

		if description != "" {
			input["description"] = description
		}

		if priority >= 0 && priority <= 4 {
			input["priority"] = priority
		}

		if assignToMe {
			viewer, err := client.GetViewer(context.Background())
			if err != nil {
				output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["assigneeId"] = viewer.ID
		}

		// Create issue
		issue, err := client.CreateIssue(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
		} else if plaintext {
			fmt.Printf("Created issue %s: %s\n", issue.Identifier, issue.Title)
		} else {
			fmt.Printf("%s Created issue %s: %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
				issue.Title)
			if issue.Assignee != nil {
				fmt.Printf("  Assigned to: %s\n", color.New(color.FgCyan).Sprint(issue.Assignee.Name))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueGetCmd)
	issueCmd.AddCommand(issueAssignCmd)
	issueCmd.AddCommand(issueCreateCmd)

	// Issue list flags
	issueListCmd.Flags().StringP("assignee", "a", "", "Filter by assignee (email or 'me')")
	issueListCmd.Flags().StringP("state", "s", "", "Filter by state name")
	issueListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	issueListCmd.Flags().IntP("priority", "r", -1, "Filter by priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueListCmd.Flags().IntP("limit", "l", 50, "Maximum number of issues to fetch")
	issueListCmd.Flags().BoolP("include-completed", "c", false, "Include completed and canceled issues")
	issueListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	issueListCmd.Flags().StringP("newer-than", "n", "", "Show issues created after this time (default: 6_months_ago, use 'all_time' for no filter)")

	// Issue create flags
	issueCreateCmd.Flags().StringP("title", "", "", "Issue title (required)")
	issueCreateCmd.Flags().StringP("description", "d", "", "Issue description")
	issueCreateCmd.Flags().StringP("team", "t", "", "Team key (required)")
	issueCreateCmd.Flags().Int("priority", 3, "Priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueCreateCmd.Flags().BoolP("assign-me", "m", false, "Assign to yourself")
	_ = issueCreateCmd.MarkFlagRequired("title")
	_ = issueCreateCmd.MarkFlagRequired("team")
}
