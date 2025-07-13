package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/pkg/auth"
	"github.com/dorkitude/linctl/pkg/output"
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

		issues, err := client.GetIssues(context.Background(), filter, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(issues.Nodes) == 0 {
			output.Info("No issues found", plaintext, jsonOut)
			return
		}

		// Prepare table data
		headers := []string{"ID", "Title", "State", "Assignee", "Team", "Priority"}
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

			rows[i] = []string{
				issue.Identifier,
				truncateString(issue.Title, 50),
				state,
				assignee,
				team,
				priority,
			}
		}

		tableData := output.TableData{
			Headers: headers,
			Rows:    rows,
		}

		output.Table(tableData, plaintext, jsonOut)

		if !plaintext && !jsonOut && issues.PageInfo.HasNextPage {
			fmt.Printf("\n%s Use --limit to see more results\n", 
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
			fmt.Printf("Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
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
			fmt.Printf("State: %s\n", 
				color.New(color.FgGreen).Sprint(issue.State.Name))
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
		fmt.Printf("Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))
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

	if state, _ := cmd.Flags().GetString("state"); state != "" {
		filter["state"] = map[string]interface{}{"name": map[string]interface{}{"eq": state}}
	}

	if team, _ := cmd.Flags().GetString("team"); team != "" {
		filter["team"] = map[string]interface{}{"key": map[string]interface{}{"eq": team}}
	}

	if priority, _ := cmd.Flags().GetInt("priority"); priority != -1 {
		filter["priority"] = map[string]interface{}{"eq": priority}
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

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueGetCmd)

	// Issue list flags
	issueListCmd.Flags().StringP("assignee", "a", "", "Filter by assignee (email or 'me')")
	issueListCmd.Flags().StringP("state", "s", "", "Filter by state name")
	issueListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	issueListCmd.Flags().IntP("priority", "r", -1, "Filter by priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueListCmd.Flags().IntP("limit", "l", 50, "Maximum number of issues to fetch")
}