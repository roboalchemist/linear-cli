package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/roboalchemist/linear-cli/pkg/utils"
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
	Long: `Manage Linear projects including listing, viewing, creating projects, and managing milestones.

Examples:
  linear-cli project list                      # List active projects
  linear-cli project list --include-completed  # List all projects including completed
  linear-cli project list --newer-than 1_month_ago  # List projects from last month
  linear-cli project get PROJECT-ID            # Get project details
  linear-cli project milestone list PROJECT-ID # List project milestones
  linear-cli project create                    # Create a new project`,
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
			// ProjectFilter uses accessibleTeams (TeamCollectionFilter) since
			// projects can belong to multiple teams. We filter for projects
			// where at least some of the accessible teams match the given team ID.
			filter["accessibleTeams"] = map[string]interface{}{
				"some": map[string]interface{}{
					"id": map[string]interface{}{"eq": team.ID},
				},
			}
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
				if project.Health != "" {
					fmt.Printf("- **Health**: %s\n", project.Health)
				}
				if project.Priority > 0 {
					fmt.Printf("- **Priority**: %s\n", project.PriorityLabel)
				}
				if project.Scope > 0 {
					fmt.Printf("- **Scope**: %.0f\n", project.Scope)
				}
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
				if project.CompletedAt != nil {
					fmt.Printf("- **Completed**: %s\n", project.CompletedAt.Format("2006-01-02"))
				}
				if project.CanceledAt != nil {
					fmt.Printf("- **Canceled**: %s\n", project.CanceledAt.Format("2006-01-02"))
				}
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
			headers := []string{"Name", "State", "Health", "Progress", "Lead", "Teams", "URL"}
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

				// Format health with color
				health := project.Health
				healthColor := color.New(color.FgWhite)
				switch project.Health {
				case "onTrack":
					healthColor = color.New(color.FgGreen)
					health = "On Track"
				case "atRisk":
					healthColor = color.New(color.FgYellow)
					health = "At Risk"
				case "offTrack":
					healthColor = color.New(color.FgRed)
					health = "Off Track"
				}

				// Format progress
				progressStr := fmt.Sprintf("%.0f%%", project.Progress*100)

				rows = append(rows, []string{
					truncateString(project.Name, 25),
					stateColor.Sprint(project.State),
					healthColor.Sprint(health),
					progressStr,
					lead,
					teams,
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
			if project.HealthUpdatedAt != nil {
				fmt.Printf("- **Health Updated**: %s\n", project.HealthUpdatedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("- **Scope**: %.0f\n", project.Scope)
			if project.Priority > 0 {
				fmt.Printf("- **Priority**: %s (%d)\n", project.PriorityLabel, project.Priority)
			}
			if project.Icon != nil && *project.Icon != "" {
				fmt.Printf("- **Icon**: %s\n", *project.Icon)
			}
			fmt.Printf("- **Color**: %s\n", project.Color)
			if project.Trashed {
				fmt.Printf("- **Trashed**: yes\n")
			}

			fmt.Printf("\n## Timeline\n")
			if project.StartDate != nil {
				dateStr := *project.StartDate
				if project.StartDateResolution != "" {
					dateStr += fmt.Sprintf(" (%s)", project.StartDateResolution)
				}
				fmt.Printf("- **Start Date**: %s\n", dateStr)
			}
			if project.TargetDate != nil {
				dateStr := *project.TargetDate
				if project.TargetDateResolution != "" {
					dateStr += fmt.Sprintf(" (%s)", project.TargetDateResolution)
				}
				fmt.Printf("- **Target Date**: %s\n", dateStr)
			}
			fmt.Printf("- **Created**: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))
			if project.StartedAt != nil {
				fmt.Printf("- **Started**: %s\n", project.StartedAt.Format("2006-01-02 15:04:05"))
			}
			if project.CompletedAt != nil {
				fmt.Printf("- **Completed**: %s\n", project.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			if project.CanceledAt != nil {
				fmt.Printf("- **Canceled**: %s\n", project.CanceledAt.Format("2006-01-02 15:04:05"))
			}
			if project.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", project.ArchivedAt.Format("2006-01-02 15:04:05"))
			}
			if project.AutoArchivedAt != nil {
				fmt.Printf("- **Auto-Archived**: %s\n", project.AutoArchivedAt.Format("2006-01-02 15:04:05"))
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

			// Milestones
			if project.ProjectMilestones != nil && len(project.ProjectMilestones.Nodes) > 0 {
				fmt.Printf("\n## Milestones\n")
				for _, ms := range project.ProjectMilestones.Nodes {
					fmt.Printf("- **%s** — %s, %.0f%%", ms.Name, ms.Status, ms.Progress*100)
					if ms.TargetDate != nil {
						fmt.Printf(", target: %s", *ms.TargetDate)
					}
					fmt.Println()
				}
			}

			// Project Updates
			if project.ProjectUpdates != nil && len(project.ProjectUpdates.Nodes) > 0 {
				fmt.Printf("\n## Recent Project Updates\n")
				for _, update := range project.ProjectUpdates.Nodes {
					fmt.Printf("\n### %s by %s\n", update.CreatedAt.Format("2006-01-02 15:04"), safeUserName(update.User))
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

			// Show milestones if available
			if project.ProjectMilestones != nil && len(project.ProjectMilestones.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Milestones:"))
				for _, ms := range project.ProjectMilestones.Nodes {
					statusColor := color.New(color.FgWhite)
					switch strings.ToLower(ms.Status) {
					case "done":
						statusColor = color.New(color.FgGreen)
					case "next":
						statusColor = color.New(color.FgBlue)
					case "overdue":
						statusColor = color.New(color.FgRed)
					}
					targetStr := ""
					if ms.TargetDate != nil {
						targetStr = fmt.Sprintf(" → %s", *ms.TargetDate)
					}
					fmt.Printf("  • %s %s %.0f%%%s\n",
						color.New(color.FgCyan).Sprint(ms.Name),
						statusColor.Sprintf("[%s]", ms.Status),
						ms.Progress*100,
						targetStr)
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

var projectAddTeamCmd = &cobra.Command{
	Use:   "add-team PROJECT-ID TEAM-KEY [TEAM-KEY...]",
	Short: "Add teams to a project",
	Long: `Add one or more teams to a project. Teams are specified by their key (e.g. ENG, DESIGN).

Examples:
  linear-cli project add-team PROJECT-ID ENG
  linear-cli project add-team PROJECT-ID ENG DESIGN OPS`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]
		teamKeys := args[1:]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		ctx := context.Background()

		// Fetch current project to get existing teams
		project, err := client.GetProject(ctx, projectID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Build set of existing team IDs
		existingIDs := make(map[string]bool)
		if project.Teams != nil {
			for _, t := range project.Teams.Nodes {
				existingIDs[t.ID] = true
			}
		}

		// Resolve each team key and collect new IDs
		var newTeamIDs []string
		for _, key := range teamKeys {
			team, err := client.GetTeam(ctx, key)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", key, err), plaintext, jsonOut)
				os.Exit(1)
			}
			if existingIDs[team.ID] {
				if !jsonOut {
					fmt.Fprintf(os.Stderr, "Team %s (%s) is already on the project, skipping\n", key, team.Name)
				}
				continue
			}
			newTeamIDs = append(newTeamIDs, team.ID)
			existingIDs[team.ID] = true
		}

		if len(newTeamIDs) == 0 {
			if jsonOut {
				output.JSON(project)
			} else {
				fmt.Println("No new teams to add.")
			}
			return
		}

		// Build full team ID list (existing + new)
		var allTeamIDs []string
		if project.Teams != nil {
			for _, t := range project.Teams.Nodes {
				allTeamIDs = append(allTeamIDs, t.ID)
			}
		}
		allTeamIDs = append(allTeamIDs, newTeamIDs...)

		// Update project
		input := map[string]interface{}{
			"teamIds": allTeamIDs,
		}
		updated, err := client.UpdateProject(ctx, projectID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(updated)
		} else if plaintext {
			fmt.Printf("# Updated Project: %s\n\n## Teams\n", updated.Name)
			if updated.Teams != nil {
				for _, t := range updated.Teams.Nodes {
					fmt.Printf("- %s (%s)\n", t.Key, t.Name)
				}
			}
		} else {
			fmt.Printf("\n%s Added %d team(s) to project %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(newTeamIDs),
				color.New(color.FgCyan, color.Bold).Sprint(updated.Name))
			if updated.Teams != nil {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Teams:"))
				for _, t := range updated.Teams.Nodes {
					fmt.Printf("  • %s - %s\n",
						color.New(color.FgCyan).Sprint(t.Key),
						t.Name)
				}
			}
			fmt.Println()
		}
	},
}

var projectRemoveTeamCmd = &cobra.Command{
	Use:   "remove-team PROJECT-ID TEAM-KEY [TEAM-KEY...]",
	Short: "Remove teams from a project",
	Long: `Remove one or more teams from a project. Teams are specified by their key (e.g. ENG, DESIGN).

Examples:
  linear-cli project remove-team PROJECT-ID ENG
  linear-cli project remove-team PROJECT-ID ENG DESIGN`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]
		teamKeys := args[1:]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		ctx := context.Background()

		// Fetch current project to get existing teams
		project, err := client.GetProject(ctx, projectID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Resolve team keys to IDs for removal
		removeIDs := make(map[string]bool)
		for _, key := range teamKeys {
			team, err := client.GetTeam(ctx, key)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", key, err), plaintext, jsonOut)
				os.Exit(1)
			}
			// Check if team is actually on the project
			found := false
			if project.Teams != nil {
				for _, t := range project.Teams.Nodes {
					if t.ID == team.ID {
						found = true
						break
					}
				}
			}
			if !found {
				if !jsonOut {
					fmt.Fprintf(os.Stderr, "Team %s (%s) is not on the project, skipping\n", key, team.Name)
				}
				continue
			}
			removeIDs[team.ID] = true
		}

		if len(removeIDs) == 0 {
			if jsonOut {
				output.JSON(project)
			} else {
				fmt.Println("No teams to remove.")
			}
			return
		}

		// Build team ID list without removed teams
		var remainingIDs []string
		if project.Teams != nil {
			for _, t := range project.Teams.Nodes {
				if !removeIDs[t.ID] {
					remainingIDs = append(remainingIDs, t.ID)
				}
			}
		}

		// Update project
		input := map[string]interface{}{
			"teamIds": remainingIDs,
		}
		updated, err := client.UpdateProject(ctx, projectID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(updated)
		} else if plaintext {
			fmt.Printf("# Updated Project: %s\n\n## Teams\n", updated.Name)
			if updated.Teams != nil {
				for _, t := range updated.Teams.Nodes {
					fmt.Printf("- %s (%s)\n", t.Key, t.Name)
				}
			}
		} else {
			fmt.Printf("\n%s Removed %d team(s) from project %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(removeIDs),
				color.New(color.FgCyan, color.Bold).Sprint(updated.Name))
			if updated.Teams != nil && len(updated.Teams.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Remaining Teams:"))
				for _, t := range updated.Teams.Nodes {
					fmt.Printf("  • %s - %s\n",
						color.New(color.FgCyan).Sprint(t.Key),
						t.Name)
				}
			} else {
				fmt.Println("  No teams remain on this project.")
			}
			fmt.Println()
		}
	},
}

var projectCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new project",
	Long: `Create a new project in Linear.

The description can be provided inline via --description or read from a markdown file via --description-file.
Use --description-file - to read from stdin.

Examples:
  linear-cli project create --name "My Project" --team-ids TEAM-UUID
  linear-cli project create --name "My Project" --description "Details" --state started
  linear-cli project create --name "My Project" --description-file project-brief.md`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		name, _ := cmd.Flags().GetString("name")
		input := map[string]interface{}{
			"name": name,
		}
		// Resolve description from --description or --description-file
		filePath, _ := cmd.Flags().GetString("description-file")
		if cmd.Flags().Changed("description") || filePath != "" {
			descFlag, _ := cmd.Flags().GetString("description")
			desc, err := resolveBodyFromFlags(descFlag, cmd.Flags().Changed("description"), filePath, "description", "description-file")
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["description"] = desc
		}
		if cmd.Flags().Changed("state") {
			s, _ := cmd.Flags().GetString("state")
			input["state"] = s
		}
		if cmd.Flags().Changed("team-ids") {
			ids, _ := cmd.Flags().GetStringSlice("team-ids")
			input["teamIds"] = ids
		}
		if cmd.Flags().Changed("start-date") {
			d, _ := cmd.Flags().GetString("start-date")
			input["startDate"] = d
		}
		if cmd.Flags().Changed("target-date") {
			d, _ := cmd.Flags().GetString("target-date")
			input["targetDate"] = d
		}

		// Handle icon
		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			input["icon"] = icon
		}

		// Handle color
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}

		// Handle priority
		if cmd.Flags().Changed("priority") {
			priority, _ := cmd.Flags().GetInt("priority")
			input["priority"] = priority
		}

		// Handle content (rich markdown)
		if cmd.Flags().Changed("content") {
			content, _ := cmd.Flags().GetString("content")
			input["content"] = content
		}

		// Handle lead
		if cmd.Flags().Changed("lead") {
			lead, _ := cmd.Flags().GetString("lead")
			switch strings.ToLower(lead) {
			case "me":
				viewer, err := client.GetViewer(context.Background())
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}
				input["leadId"] = viewer.ID
			case "":
				// Don't set leadId
			default:
				users, err := client.GetUsers(context.Background(), 100, "", "")
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get users: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}
				var foundUser *api.User
				for _, user := range users.Nodes {
					if user.ID == lead || user.Email == lead || user.Name == lead {
						foundUser = &user
						break
					}
				}
				if foundUser == nil {
					output.Error(fmt.Sprintf("User not found: %s", lead), plaintext, jsonOut)
					os.Exit(1)
				}
				input["leadId"] = foundUser.ID
			}
		}

		// Handle members
		if cmd.Flags().Changed("members") {
			membersArg, _ := cmd.Flags().GetStringSlice("members")
			if len(membersArg) > 0 {
				users, err := client.GetUsers(context.Background(), 100, "", "")
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get users: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}
				var memberIDs []string
				for _, member := range membersArg {
					memberLower := strings.ToLower(member)
					if memberLower == "me" {
						viewer, err := client.GetViewer(context.Background())
						if err != nil {
							output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
							os.Exit(1)
						}
						memberIDs = append(memberIDs, viewer.ID)
						continue
					}
					var foundUser *api.User
					for _, user := range users.Nodes {
						if strings.EqualFold(user.Email, member) || strings.EqualFold(user.Name, member) {
							foundUser = &user
							break
						}
					}
					if foundUser == nil {
						output.Error(fmt.Sprintf("User not found: %s", member), plaintext, jsonOut)
						os.Exit(1)
					}
					memberIDs = append(memberIDs, foundUser.ID)
				}
				input["memberIds"] = memberIDs
			}
		}

		// Handle template
		if cmd.Flags().Changed("template-id") {
			templateID, _ := cmd.Flags().GetString("template-id")
			input["templateId"] = templateID
		}

		// Handle use-default-template
		if cmd.Flags().Changed("use-default-template") {
			useDefault, _ := cmd.Flags().GetBool("use-default-template")
			input["useDefaultTemplate"] = useDefault
		}

		// Handle converted-from-issue
		if cmd.Flags().Changed("converted-from-issue") {
			issueID, _ := cmd.Flags().GetString("converted-from-issue")
			input["convertedFromIssueId"] = issueID
		}

		// Handle date resolution
		if cmd.Flags().Changed("start-date-resolution") {
			res, _ := cmd.Flags().GetString("start-date-resolution")
			input["startDateResolution"] = res
		}
		if cmd.Flags().Changed("target-date-resolution") {
			res, _ := cmd.Flags().GetString("target-date-resolution")
			input["targetDateResolution"] = res
		}

		project, err := client.CreateProject(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(project)
		} else {
			output.Success(fmt.Sprintf("Created project %s (%s)",
				color.New(color.FgWhite, color.Bold).Sprint(project.Name),
				project.State), plaintext, jsonOut)
		}
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:     "update PROJECT-ID",
	Aliases: []string{"edit"},
	Short:   "Update a project",
	Long: `Update a project's name, description, state, dates, or lead.

The description can be provided inline via --description or read from a markdown file via --description-file.
Use --description-file - to read from stdin.

Examples:
  linear-cli project update PROJECT-ID --name "New Name"
  linear-cli project update PROJECT-ID --state started
  linear-cli project update PROJECT-ID --lead user@example.com
  linear-cli project update PROJECT-ID --lead me
  linear-cli project update PROJECT-ID --lead none`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			n, _ := cmd.Flags().GetString("name")
			input["name"] = n
		}
		// Resolve description from --description or --description-file
		filePath, _ := cmd.Flags().GetString("description-file")
		if cmd.Flags().Changed("description") || filePath != "" {
			descFlag, _ := cmd.Flags().GetString("description")
			desc, err := resolveBodyFromFlags(descFlag, cmd.Flags().Changed("description"), filePath, "description", "description-file")
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["description"] = desc
		}
		if cmd.Flags().Changed("state") {
			s, _ := cmd.Flags().GetString("state")
			input["state"] = s
		}
		if cmd.Flags().Changed("start-date") {
			d, _ := cmd.Flags().GetString("start-date")
			input["startDate"] = d
		}
		if cmd.Flags().Changed("target-date") {
			d, _ := cmd.Flags().GetString("target-date")
			input["targetDate"] = d
		}

		// Handle lead update
		if cmd.Flags().Changed("lead") {
			lead, _ := cmd.Flags().GetString("lead")
			switch strings.ToLower(lead) {
			case "me":
				// Get current user
				viewer, err := client.GetViewer(context.Background())
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}
				input["leadId"] = viewer.ID
			case "none", "unassigned", "":
				input["leadId"] = nil
			default:
				// Look up user by email or name
				users, err := client.GetUsers(context.Background(), 100, "", "")
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get users: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}

				var foundUser *api.User
				for _, user := range users.Nodes {
					if user.ID == lead || user.Email == lead || user.Name == lead {
						foundUser = &user
						break
					}
				}

				if foundUser == nil {
					output.Error(fmt.Sprintf("User not found: %s", lead), plaintext, jsonOut)
					os.Exit(1)
				}

				input["leadId"] = foundUser.ID
			}
		}

		// Handle icon update
		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			if icon == "" {
				input["icon"] = nil
			} else {
				input["icon"] = icon
			}
		}

		// Handle color update
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}

		// Handle priority update
		if cmd.Flags().Changed("priority") {
			priority, _ := cmd.Flags().GetInt("priority")
			input["priority"] = priority
		}

		// Handle content update (rich markdown content)
		if cmd.Flags().Changed("content") {
			content, _ := cmd.Flags().GetString("content")
			input["content"] = content
		}

		// Handle members update
		if cmd.Flags().Changed("members") {
			membersArg, _ := cmd.Flags().GetStringSlice("members")
			if len(membersArg) == 1 && strings.ToLower(membersArg[0]) == "none" {
				input["memberIds"] = []string{}
			} else {
				// Resolve member emails/names to IDs
				users, err := client.GetUsers(context.Background(), 100, "", "")
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get users: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}

				var memberIDs []string
				for _, member := range membersArg {
					var foundUser *api.User
					memberLower := strings.ToLower(member)
					if memberLower == "me" {
						viewer, err := client.GetViewer(context.Background())
						if err != nil {
							output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
							os.Exit(1)
						}
						memberIDs = append(memberIDs, viewer.ID)
						continue
					}
					for _, user := range users.Nodes {
						if strings.EqualFold(user.Email, member) || strings.EqualFold(user.Name, member) {
							foundUser = &user
							break
						}
					}
					if foundUser == nil {
						output.Error(fmt.Sprintf("User not found: %s", member), plaintext, jsonOut)
						os.Exit(1)
					}
					memberIDs = append(memberIDs, foundUser.ID)
				}
				input["memberIds"] = memberIDs
			}
		}

		// Handle Slack integration flags
		if cmd.Flags().Changed("slack-new-issue") {
			v, _ := cmd.Flags().GetBool("slack-new-issue")
			input["slackNewIssue"] = v
		}
		if cmd.Flags().Changed("slack-issue-comments") {
			v, _ := cmd.Flags().GetBool("slack-issue-comments")
			input["slackIssueComments"] = v
		}
		if cmd.Flags().Changed("slack-issue-statuses") {
			v, _ := cmd.Flags().GetBool("slack-issue-statuses")
			input["slackIssueStatuses"] = v
		}

		// Handle trashed
		if cmd.Flags().Changed("trashed") {
			v, _ := cmd.Flags().GetBool("trashed")
			input["trashed"] = v
		}

		// Handle completed-at and canceled-at
		if cmd.Flags().Changed("completed-at") {
			v, _ := cmd.Flags().GetString("completed-at")
			if v == "" || strings.ToLower(v) == "none" {
				input["completedAt"] = nil
			} else {
				input["completedAt"] = v
			}
		}
		if cmd.Flags().Changed("canceled-at") {
			v, _ := cmd.Flags().GetString("canceled-at")
			if v == "" || strings.ToLower(v) == "none" {
				input["canceledAt"] = nil
			} else {
				input["canceledAt"] = v
			}
		}

		// Handle converted-from-issue
		if cmd.Flags().Changed("converted-from-issue") {
			v, _ := cmd.Flags().GetString("converted-from-issue")
			if v == "" || strings.ToLower(v) == "none" {
				input["convertedFromIssueId"] = nil
			} else {
				input["convertedFromIssueId"] = v
			}
		}

		// Handle last-applied-template
		if cmd.Flags().Changed("last-applied-template") {
			v, _ := cmd.Flags().GetString("last-applied-template")
			if v == "" || strings.ToLower(v) == "none" {
				input["lastAppliedTemplateId"] = nil
			} else {
				input["lastAppliedTemplateId"] = v
			}
		}

		// Handle date resolution
		if cmd.Flags().Changed("start-date-resolution") {
			v, _ := cmd.Flags().GetString("start-date-resolution")
			input["startDateResolution"] = v
		}
		if cmd.Flags().Changed("target-date-resolution") {
			v, _ := cmd.Flags().GetString("target-date-resolution")
			input["targetDateResolution"] = v
		}

		// Handle update reminder settings
		if cmd.Flags().Changed("update-reminder-frequency") {
			v, _ := cmd.Flags().GetFloat64("update-reminder-frequency")
			input["updateReminderFrequency"] = v
		}
		if cmd.Flags().Changed("frequency-resolution") {
			v, _ := cmd.Flags().GetString("frequency-resolution")
			input["frequencyResolution"] = v
		}
		if cmd.Flags().Changed("update-reminders-day") {
			v, _ := cmd.Flags().GetString("update-reminders-day")
			input["updateRemindersDay"] = v
		}
		if cmd.Flags().Changed("update-reminders-hour") {
			v, _ := cmd.Flags().GetInt("update-reminders-hour")
			input["updateRemindersHour"] = v
		}
		if cmd.Flags().Changed("update-reminders-paused-until") {
			v, _ := cmd.Flags().GetString("update-reminders-paused-until")
			if v == "" || strings.ToLower(v) == "none" {
				input["projectUpdateRemindersPausedUntilAt"] = nil
			} else {
				input["projectUpdateRemindersPausedUntilAt"] = v
			}
		}

		if len(input) == 0 {
			output.Error("No fields to update.", plaintext, jsonOut)
			os.Exit(1)
		}

		project, err := client.UpdateProject(context.Background(), projectID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(project)
		} else {
			leadInfo := ""
			if project.Lead != nil {
				leadInfo = fmt.Sprintf(" (lead: %s)", project.Lead.Name)
			}
			output.Success(fmt.Sprintf("Updated project %s%s",
				color.New(color.FgWhite, color.Bold).Sprint(project.Name), leadInfo), plaintext, jsonOut)
		}
	},
}

var projectArchiveCmd = &cobra.Command{
	Use:   "archive PROJECT-ID",
	Short: "Archive a project",
	Long:  `Archive a project. Archived projects can be restored in the Linear UI.`,
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
		err = client.ArchiveProject(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to archive project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Archived project", plaintext, jsonOut)
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete PROJECT-ID",
	Short: "Permanently delete a project",
	Long:  `Permanently delete a project. This cannot be undone.`,
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
		err = client.DeleteProject(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted project", plaintext, jsonOut)
	},
}

var projectIssuesCmd = &cobra.Command{
	Use:     "issues PROJECT-ID",
	Aliases: []string{"issue"},
	Short:   "List issues in a project",
	Long: `List all issues that belong to a specific project.

Examples:
  linear-cli project issues PROJECT-ID            # List project issues
  linear-cli project issues PROJECT-ID --json     # JSON output`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")

		issues, err := client.GetProjectIssues(context.Background(), projectID, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get project issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issues.Nodes)
			return
		}

		if len(issues.Nodes) == 0 {
			if plaintext {
				fmt.Println("No issues found")
			} else {
				fmt.Printf("\n%s No issues in this project\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Issues")
			fmt.Println("ID\tTitle\tState\tPriority\tAssignee")
			for _, i := range issues.Nodes {
				state := ""
				if i.State != nil {
					state = i.State.Name
				}
				assignee := "Unassigned"
				if i.Assignee != nil {
					assignee = i.Assignee.Name
				}
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
					i.Identifier, i.Title, state, i.PriorityLabel, assignee)
			}
		} else {
			headers := []string{"ID", "Title", "State", "Priority", "Assignee"}
			rows := [][]string{}

			for _, i := range issues.Nodes {
				state := ""
				stateColor := color.New(color.FgWhite)
				if i.State != nil {
					state = i.State.Name
					switch i.State.Type {
					case "triage":
						stateColor = color.New(color.FgMagenta)
					case "backlog":
						stateColor = color.New(color.FgCyan)
					case "started":
						stateColor = color.New(color.FgYellow)
					case "completed":
						stateColor = color.New(color.FgGreen)
					case "canceled":
						stateColor = color.New(color.FgRed)
					}
				}
				assignee := "Unassigned"
				if i.Assignee != nil {
					assignee = i.Assignee.Name
				}

				rows = append(rows, []string{
					color.New(color.FgCyan).Sprint(i.Identifier),
					i.Title,
					stateColor.Sprint(state),
					i.PriorityLabel,
					assignee,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d issues in project\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(issues.Nodes))
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectAddTeamCmd)
	projectCmd.AddCommand(projectRemoveTeamCmd)
	projectCmd.AddCommand(projectIssuesCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectArchiveCmd)
	projectCmd.AddCommand(projectDeleteCmd)

	// Project issues flags
	projectIssuesCmd.Flags().IntP("limit", "l", 50, "Maximum number of issues to return")

	// Project create flags
	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().StringP("description", "d", "", "Project description")
	projectCreateCmd.Flags().String("description-file", "", "Read description from a markdown file (use - for stdin)")
	projectCreateCmd.Flags().StringSlice("team-ids", nil, "Team IDs to associate with")
	projectCreateCmd.Flags().String("state", "planned", "State: planned, started, paused, completed, canceled")
	projectCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	projectCreateCmd.Flags().String("target-date", "", "Target date (YYYY-MM-DD)")
	projectCreateCmd.Flags().String("icon", "", "Project icon (emoji)")
	projectCreateCmd.Flags().StringP("color", "c", "", "Project color (hex, e.g., #4285F4)")
	projectCreateCmd.Flags().Int("priority", -1, "Project priority (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)")
	projectCreateCmd.Flags().String("content", "", "Project content (rich markdown)")
	projectCreateCmd.Flags().StringP("lead", "L", "", "Project lead (email, name, UUID, or 'me')")
	projectCreateCmd.Flags().StringSlice("members", nil, "Project members (emails/names, repeatable)")
	projectCreateCmd.Flags().String("template-id", "", "Template ID to apply")
	projectCreateCmd.Flags().Bool("use-default-template", false, "Apply default project template")
	projectCreateCmd.Flags().String("converted-from-issue", "", "Issue ID this project was converted from")
	projectCreateCmd.Flags().String("start-date-resolution", "", "Start date resolution (day, month, quarter, year)")
	projectCreateCmd.Flags().String("target-date-resolution", "", "Target date resolution (day, month, quarter, year)")
	_ = projectCreateCmd.MarkFlagRequired("name")

	// Project update flags
	projectUpdateCmd.Flags().String("name", "", "New project name")
	projectUpdateCmd.Flags().StringP("description", "d", "", "New description")
	projectUpdateCmd.Flags().String("description-file", "", "Read description from a markdown file (use - for stdin)")
	projectUpdateCmd.Flags().String("state", "", "New state: planned, started, paused, completed, canceled")
	projectUpdateCmd.Flags().String("start-date", "", "New start date (YYYY-MM-DD)")
	projectUpdateCmd.Flags().String("target-date", "", "New target date (YYYY-MM-DD)")
	projectUpdateCmd.Flags().StringP("lead", "L", "", "Project lead (email, name, UUID, 'me', or 'none' to unset)")
	projectUpdateCmd.Flags().String("icon", "", "Project icon (emoji or empty to remove)")
	projectUpdateCmd.Flags().StringP("color", "c", "", "Project color (hex, e.g., #4285F4)")
	projectUpdateCmd.Flags().Int("priority", -1, "Project priority (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)")
	projectUpdateCmd.Flags().String("content", "", "Project content (rich markdown)")
	projectUpdateCmd.Flags().StringSlice("members", nil, "Project members (emails/names, repeatable, or 'none' to clear)")
	// Slack integration flags
	projectUpdateCmd.Flags().Bool("slack-new-issue", false, "Notify Slack on new issues")
	projectUpdateCmd.Flags().Bool("slack-issue-comments", false, "Notify Slack on issue comments")
	projectUpdateCmd.Flags().Bool("slack-issue-statuses", false, "Notify Slack on issue status changes")
	// Additional update flags
	projectUpdateCmd.Flags().Bool("trashed", false, "Mark project as trashed")
	projectUpdateCmd.Flags().String("completed-at", "", "Completion timestamp (ISO 8601 or 'none' to unset)")
	projectUpdateCmd.Flags().String("canceled-at", "", "Cancellation timestamp (ISO 8601 or 'none' to unset)")
	projectUpdateCmd.Flags().String("converted-from-issue", "", "Issue ID this project was converted from (or 'none' to unset)")
	projectUpdateCmd.Flags().String("last-applied-template", "", "Last applied template ID (or 'none' to unset)")
	// Date resolution flags
	projectUpdateCmd.Flags().String("start-date-resolution", "", "Start date resolution (day, month, quarter, year)")
	projectUpdateCmd.Flags().String("target-date-resolution", "", "Target date resolution (day, month, quarter, year)")
	// Update reminder flags
	projectUpdateCmd.Flags().Float64("update-reminder-frequency", 0, "Reminder frequency (number of periods)")
	projectUpdateCmd.Flags().String("frequency-resolution", "", "Frequency resolution (day, week, month)")
	projectUpdateCmd.Flags().String("update-reminders-day", "", "Day for update reminders (Monday, Tuesday, etc.)")
	projectUpdateCmd.Flags().Int("update-reminders-hour", -1, "Hour for update reminders (0-23)")
	projectUpdateCmd.Flags().String("update-reminders-paused-until", "", "Pause reminders until (ISO 8601 timestamp or 'none' to resume)")

	// List command flags
	projectListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	projectListCmd.Flags().StringP("state", "s", "", "Filter by state (planned, started, paused, completed, canceled)")
	projectListCmd.Flags().IntP("limit", "l", 50, "Maximum number of projects to return")
	projectListCmd.Flags().BoolP("include-completed", "c", false, "Include completed and canceled projects")
	projectListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	projectListCmd.Flags().StringP("newer-than", "n", "", "Show projects created after this time (default: 6_months_ago, use 'all_time' for no filter)")
}
