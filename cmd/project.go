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

// constructProjectURL constructs an ID-based project URL
func constructProjectURL(projectID string, originalURL string) string {
	// Extract workspace from the original URL
	// Format: https://linear.app/{workspace}/project/{slug}
	if originalURL == "" {
		return ""
	}
	
	parts := strings.Split(originalURL, "/")
	if len(parts) >= 5 {
		workspace := parts[3]
		return fmt.Sprintf("https://linear.app/%s/project/%s", workspace, projectID)
	}
	
	// Fallback to original URL if we can't parse it
	return originalURL
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Linear projects",
	Long: `Manage Linear projects including listing, viewing, and creating projects.

Examples:
  linctl project list                      # List active projects
  linctl project list --include-completed  # List all projects including completed
  linctl project list --newer-than 1_month_ago  # List projects from last month
  linctl project get PROJECT-ID            # Get project details
  linctl project create                    # Create a new project`,
}

var projectListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects",
	Long:    `List all projects in your Linear workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		// Get filters
		teamKey, _ := cmd.Flags().GetString("team")
		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		includeCompleted, _ := cmd.Flags().GetBool("include-completed")

		// Build filter
		filter := make(map[string]interface{})
		if teamKey != "" {
			// Get team ID from key
			team, err := client.GetTeam(context.Background(), teamKey)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
				os.Exit(1)
			}
			filter["team"] = map[string]interface{}{"id": team.ID}
		}
		if state != "" {
			filter["state"] = map[string]interface{}{"eq": state}
		} else if !includeCompleted {
			// Only filter out completed projects if no specific state is requested
			filter["state"] = map[string]interface{}{
				"nin": []string{"completed", "canceled"},
			}
		}

		// Handle newer-than filter
		newerThan, _ := cmd.Flags().GetString("newer-than")
		createdAt, err := utils.ParseTimeExpression(newerThan)
		if err != nil {
			output.Error(fmt.Sprintf("Invalid newer-than value: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if createdAt != "" {
			filter["createdAt"] = map[string]interface{}{"gte": createdAt}
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

		// Get projects
		projects, err := client.GetProjects(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list projects: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(projects.Nodes)
			return
		} else if plaintext {
			fmt.Println("# Projects")
			for _, project := range projects.Nodes {
				fmt.Printf("## %s\n", project.Name)
				fmt.Printf("- **ID**: %s\n", project.ID)
				fmt.Printf("- **State**: %s\n", project.State)
				fmt.Printf("- **Progress**: %.0f%%\n", project.Progress*100)
				if project.Lead != nil {
					fmt.Printf("- **Lead**: %s\n", project.Lead.Name)
				} else {
					fmt.Printf("- **Lead**: Unassigned\n")
				}
				if project.Teams != nil && len(project.Teams.Nodes) > 0 {
					teams := ""
					for i, team := range project.Teams.Nodes {
						if i > 0 {
							teams += ", "
						}
						teams += team.Key
					}
					fmt.Printf("- **Teams**: %s\n", teams)
				}
				if project.StartDate != nil {
					fmt.Printf("- **Start Date**: %s\n", *project.StartDate)
				}
				if project.TargetDate != nil {
					fmt.Printf("- **Target Date**: %s\n", *project.TargetDate)
				}
				fmt.Printf("- **Created**: %s\n", project.CreatedAt.Format("2006-01-02"))
				fmt.Printf("- **Updated**: %s\n", project.UpdatedAt.Format("2006-01-02"))
				fmt.Printf("- **URL**: %s\n", constructProjectURL(project.ID, project.URL))
				if project.Description != "" {
					fmt.Printf("- **Description**: %s\n", project.Description)
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d projects\n", len(projects.Nodes))
			return
		} else {
			// Table output
			headers := []string{"Name", "State", "Lead", "Teams", "Created", "Updated", "URL"}
			rows := [][]string{}

			for _, project := range projects.Nodes {
				lead := color.New(color.FgYellow).Sprint("Unassigned")
				if project.Lead != nil {
					lead = project.Lead.Name
				}

				teams := ""
				if project.Teams != nil && len(project.Teams.Nodes) > 0 {
					for i, team := range project.Teams.Nodes {
						if i > 0 {
							teams += ", "
						}
						teams += team.Key
					}
				}

				stateColor := color.New(color.FgGreen)
				switch project.State {
				case "planned":
					stateColor = color.New(color.FgCyan)
				case "started":
					stateColor = color.New(color.FgBlue)
				case "paused":
					stateColor = color.New(color.FgYellow)
				case "completed":
					stateColor = color.New(color.FgGreen)
				case "canceled":
					stateColor = color.New(color.FgRed)
				}

				rows = append(rows, []string{
					truncateString(project.Name, 25),
					stateColor.Sprint(project.State),
					lead,
					teams,
					project.CreatedAt.Format("2006-01-02"),
					project.UpdatedAt.Format("2006-01-02"),
					constructProjectURL(project.ID, project.URL),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			if !plaintext && !jsonOut {
				fmt.Printf("\n%s %d projects\n",
					color.New(color.FgGreen).Sprint("✓"),
					len(projects.Nodes))
			}
		}
	},
}

var projectGetCmd = &cobra.Command{
	Use:     "get PROJECT-ID",
	Aliases: []string{"show"},
	Short:   "Get project details",
	Long:    `Get detailed information about a specific project.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get project details
		project, err := client.GetProject(context.Background(), projectID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(project)
		} else if plaintext {
			fmt.Printf("# %s\n\n", project.Name)
			
			if project.Description != "" {
				fmt.Printf("## Description\n%s\n\n", project.Description)
			}
			
			if project.Content != "" {
				fmt.Printf("## Content\n%s\n\n", project.Content)
			}
			
			fmt.Printf("## Core Details\n")
			fmt.Printf("- **ID**: %s\n", project.ID)
			fmt.Printf("- **Slug ID**: %s\n", project.SlugId)
			fmt.Printf("- **State**: %s\n", project.State)
			fmt.Printf("- **Progress**: %.0f%%\n", project.Progress*100)
			fmt.Printf("- **Health**: %s\n", project.Health)
			fmt.Printf("- **Scope**: %d\n", project.Scope)
			if project.Icon != nil && *project.Icon != "" {
				fmt.Printf("- **Icon**: %s\n", *project.Icon)
			}
			fmt.Printf("- **Color**: %s\n", project.Color)
			
			fmt.Printf("\n## Timeline\n")
			if project.StartDate != nil {
				fmt.Printf("- **Start Date**: %s\n", *project.StartDate)
			}
			if project.TargetDate != nil {
				fmt.Printf("- **Target Date**: %s\n", *project.TargetDate)
			}
			fmt.Printf("- **Created**: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))
			if project.CompletedAt != nil {
				fmt.Printf("- **Completed**: %s\n", project.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			if project.CanceledAt != nil {
				fmt.Printf("- **Canceled**: %s\n", project.CanceledAt.Format("2006-01-02 15:04:05"))
			}
			if project.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", project.ArchivedAt.Format("2006-01-02 15:04:05"))
			}
			
			fmt.Printf("\n## People\n")
			if project.Lead != nil {
				fmt.Printf("- **Lead**: %s (%s)\n", project.Lead.Name, project.Lead.Email)
				if project.Lead.DisplayName != "" && project.Lead.DisplayName != project.Lead.Name {
					fmt.Printf("  - Display Name: %s\n", project.Lead.DisplayName)
				}
			} else {
				fmt.Printf("- **Lead**: Unassigned\n")
			}
			if project.Creator != nil {
				fmt.Printf("- **Creator**: %s (%s)\n", project.Creator.Name, project.Creator.Email)
			}
			
			fmt.Printf("\n## Slack Integration\n")
			fmt.Printf("- **Slack New Issue**: %v\n", project.SlackNewIssue)
			fmt.Printf("- **Slack Issue Comments**: %v\n", project.SlackIssueComments)
			fmt.Printf("- **Slack Issue Statuses**: %v\n", project.SlackIssueStatuses)
			
			if project.ConvertedFromIssue != nil {
				fmt.Printf("\n## Origin\n")
				fmt.Printf("- **Converted from Issue**: %s - %s\n", project.ConvertedFromIssue.Identifier, project.ConvertedFromIssue.Title)
			}
			
			if project.LastAppliedTemplate != nil {
				fmt.Printf("\n## Template\n")
				fmt.Printf("- **Last Applied**: %s\n", project.LastAppliedTemplate.Name)
				if project.LastAppliedTemplate.Description != "" {
					fmt.Printf("  - Description: %s\n", project.LastAppliedTemplate.Description)
				}
			}
			
			
			// Teams
			if project.Teams != nil && len(project.Teams.Nodes) > 0 {
				fmt.Printf("\n## Teams\n")
				for _, team := range project.Teams.Nodes {
					fmt.Printf("- **%s** (%s)\n", team.Name, team.Key)
					if team.Description != "" {
						fmt.Printf("  - Description: %s\n", team.Description)
					}
					fmt.Printf("  - Cycles Enabled: %v\n", team.CyclesEnabled)
				}
			}
			
			fmt.Printf("\n## URL\n")
			fmt.Printf("- %s\n", constructProjectURL(project.ID, project.URL))
			
			// Show members if available
			if project.Members != nil && len(project.Members.Nodes) > 0 {
				fmt.Printf("\n## Members\n")
				for _, member := range project.Members.Nodes {
					fmt.Printf("- %s (%s)", member.Name, member.Email)
					if member.DisplayName != "" && member.DisplayName != member.Name {
						fmt.Printf(" - %s", member.DisplayName)
					}
					if member.Admin {
						fmt.Printf(" [Admin]")
					}
					if !member.Active {
						fmt.Printf(" [Inactive]")
					}
					fmt.Println()
				}
			}
			
			// Project Updates
			if project.ProjectUpdates != nil && len(project.ProjectUpdates.Nodes) > 0 {
				fmt.Printf("\n## Recent Project Updates\n")
				for _, update := range project.ProjectUpdates.Nodes {
					fmt.Printf("\n### %s by %s\n", update.CreatedAt.Format("2006-01-02 15:04"), update.User.Name)
					if update.EditedAt != nil {
						fmt.Printf("*(edited %s)*\n", update.EditedAt.Format("2006-01-02 15:04"))
					}
					fmt.Printf("- **Health**: %s\n", update.Health)
					fmt.Printf("\n%s\n", update.Body)
				}
			}
			
			// Documents
			if project.Documents != nil && len(project.Documents.Nodes) > 0 {
				fmt.Printf("\n## Documents\n")
				for _, doc := range project.Documents.Nodes {
					fmt.Printf("\n### %s\n", doc.Title)
					if doc.Icon != nil && *doc.Icon != "" {
						fmt.Printf("- **Icon**: %s\n", *doc.Icon)
					}
					fmt.Printf("- **Color**: %s\n", doc.Color)
					fmt.Printf("- **Created**: %s by %s\n", doc.CreatedAt.Format("2006-01-02"), doc.Creator.Name)
					if doc.UpdatedBy != nil {
						fmt.Printf("- **Updated**: %s by %s\n", doc.UpdatedAt.Format("2006-01-02"), doc.UpdatedBy.Name)
					}
					fmt.Printf("\n%s\n", doc.Content)
				}
			}
			
			
			// Show recent issues
			if project.Issues != nil && len(project.Issues.Nodes) > 0 {
				fmt.Printf("\n## Issues (%d total)\n", len(project.Issues.Nodes))
				for _, issue := range project.Issues.Nodes {
					stateStr := ""
					if issue.State != nil {
						switch issue.State.Type {
						case "completed":
							stateStr = "[x]"
						case "started":
							stateStr = "[~]"
						case "canceled":
							stateStr = "[-]"
						default:
							stateStr = "[ ]"
						}
					} else {
						stateStr = "[ ]"
					}
					
					assignee := "Unassigned"
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}
					
					fmt.Printf("\n### %s %s (#%d)\n", stateStr, issue.Identifier, issue.Number)
					fmt.Printf("**%s**\n", issue.Title)
					fmt.Printf("- Assignee: %s\n", assignee)
					fmt.Printf("- Priority: %s\n", priorityToString(issue.Priority))
					if issue.Estimate != nil {
						fmt.Printf("- Estimate: %.1f\n", *issue.Estimate)
					}
					if issue.State != nil {
						fmt.Printf("- State: %s\n", issue.State.Name)
					}
					if issue.Labels != nil && len(issue.Labels.Nodes) > 0 {
						labels := []string{}
						for _, label := range issue.Labels.Nodes {
							labels = append(labels, label.Name)
						}
						fmt.Printf("- Labels: %s\n", strings.Join(labels, ", "))
					}
					fmt.Printf("- Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04"))
					if issue.Description != "" {
						// Show first 3 lines of description
						lines := strings.Split(issue.Description, "\n")
						preview := ""
						for i, line := range lines {
							if i >= 3 {
								preview += "\n  ..."
								break
							}
							if i > 0 {
								preview += "\n  "
							}
							preview += line
						}
						fmt.Printf("- Description: %s\n", preview)
					}
				}
			}
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s\n", color.New(color.FgCyan, color.Bold).Sprint("📁 Project:"), project.Name)
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), project.ID)

			if project.Description != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Description:"), project.Description)
			}

			stateColor := color.New(color.FgGreen)
			switch project.State {
			case "planned":
				stateColor = color.New(color.FgCyan)
			case "started":
				stateColor = color.New(color.FgBlue)
			case "paused":
				stateColor = color.New(color.FgYellow)
			case "completed":
				stateColor = color.New(color.FgGreen)
			case "canceled":
				stateColor = color.New(color.FgRed)
			}
			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("State:"), stateColor.Sprint(project.State))

			progressColor := color.New(color.FgRed)
			if project.Progress >= 0.75 {
				progressColor = color.New(color.FgGreen)
			} else if project.Progress >= 0.5 {
				progressColor = color.New(color.FgYellow)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Progress:"), progressColor.Sprintf("%.0f%%", project.Progress*100))

			if project.StartDate != nil || project.TargetDate != nil {
				fmt.Println()
				if project.StartDate != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Start Date:"), *project.StartDate)
				}
				if project.TargetDate != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Target Date:"), *project.TargetDate)
				}
			}

			if project.Lead != nil {
				fmt.Printf("\n%s %s (%s)\n",
					color.New(color.Bold).Sprint("Lead:"),
					project.Lead.Name,
					color.New(color.FgCyan).Sprint(project.Lead.Email))
			}

			if project.Teams != nil && len(project.Teams.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Teams:"))
				for _, team := range project.Teams.Nodes {
					fmt.Printf("  • %s - %s\n",
						color.New(color.FgCyan).Sprint(team.Key),
						team.Name)
				}
			}

			// Show members if available
			if project.Members != nil && len(project.Members.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Members:"))
				for _, member := range project.Members.Nodes {
					fmt.Printf("  • %s (%s)\n",
						member.Name,
						color.New(color.FgCyan).Sprint(member.Email))
				}
			}



			// Show sample issues if available
			if project.Issues != nil && len(project.Issues.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Recent Issues:"))
				for i, issue := range project.Issues.Nodes {
					if i >= 5 {
						break // Show only first 5
					}
					stateIcon := "○"
					if issue.State != nil {
						switch issue.State.Type {
						case "completed":
							stateIcon = color.New(color.FgGreen).Sprint("✓")
						case "started":
							stateIcon = color.New(color.FgBlue).Sprint("◐")
						case "canceled":
							stateIcon = color.New(color.FgRed).Sprint("✗")
						}
					}
					assignee := "Unassigned"
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}
					fmt.Printf("  %s %s %s (%s)\n",
						stateIcon,
						color.New(color.FgCyan).Sprint(issue.Identifier),
						issue.Title,
						color.New(color.FgWhite, color.Faint).Sprint(assignee))
				}
			}

			// Show timestamps
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Timeline:"))
			fmt.Printf("  Created: %s\n", project.CreatedAt.Format("2006-01-02"))
			fmt.Printf("  Updated: %s\n", project.UpdatedAt.Format("2006-01-02"))
			if project.CompletedAt != nil {
				fmt.Printf("  Completed: %s\n", project.CompletedAt.Format("2006-01-02"))
			}
			if project.CanceledAt != nil {
				fmt.Printf("  Canceled: %s\n", project.CanceledAt.Format("2006-01-02"))
			}

			// Show URL
			if project.URL != "" {
				fmt.Printf("\n%s %s\n",
					color.New(color.Bold).Sprint("URL:"),
					color.New(color.FgBlue, color.Underline).Sprint(constructProjectURL(project.ID, project.URL)))
			}

			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)

	// List command flags
	projectListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	projectListCmd.Flags().StringP("state", "s", "", "Filter by state (planned, started, paused, completed, canceled)")
	projectListCmd.Flags().IntP("limit", "l", 50, "Maximum number of projects to return")
	projectListCmd.Flags().BoolP("include-completed", "c", false, "Include completed and canceled projects")
	projectListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	projectListCmd.Flags().StringP("newer-than", "n", "", "Show projects created after this time (default: 6_months_ago, use 'all_time' for no filter)")
}
