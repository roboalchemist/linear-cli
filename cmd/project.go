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
			fmt.Println("ID\tName\tState\tProgress\tLead\tTeams\tCreated\tUpdated")
			for _, project := range projects.Nodes {
				lead := "Unassigned"
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
				fmt.Printf("%s\t%s\t%s\t%.0f%%\t%s\t%s\t%s\t%s\n",
					project.ID,
					truncateString(project.Name, 30),
					project.State,
					project.Progress*100,
					lead,
					teams,
					project.CreatedAt.Format("2006-01-02"),
					project.UpdatedAt.Format("2006-01-02"),
				)
			}
		} else {
			// Table output
			headers := []string{"ID", "Name", "State", "Progress", "Lead", "Teams", "Created", "Updated"}
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

				progressStr := fmt.Sprintf("%.0f%%", project.Progress*100)
				progressColor := color.New(color.FgRed)
				if project.Progress >= 0.75 {
					progressColor = color.New(color.FgGreen)
				} else if project.Progress >= 0.5 {
					progressColor = color.New(color.FgYellow)
				}

				rows = append(rows, []string{
					project.ID,
					truncateString(project.Name, 25),
					stateColor.Sprint(project.State),
					progressColor.Sprint(progressStr),
					lead,
					teams,
					project.CreatedAt.Format("2006-01-02"),
					project.UpdatedAt.Format("2006-01-02"),
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
			fmt.Printf("ID: %s\n", project.ID)
			fmt.Printf("Name: %s\n", project.Name)
			if project.Description != "" {
				fmt.Printf("Description: %s\n", project.Description)
			}
			fmt.Printf("State: %s\n", project.State)
			fmt.Printf("Progress: %.0f%%\n", project.Progress*100)
			if project.StartDate != nil {
				fmt.Printf("Start Date: %s\n", *project.StartDate)
			}
			if project.TargetDate != nil {
				fmt.Printf("Target Date: %s\n", *project.TargetDate)
			}
			if project.Lead != nil {
				fmt.Printf("Lead: %s (%s)\n", project.Lead.Name, project.Lead.Email)
			}
			if project.Creator != nil {
				fmt.Printf("Creator: %s\n", project.Creator.Name)
			}
			if project.Teams != nil && len(project.Teams.Nodes) > 0 {
				fmt.Printf("Teams: ")
				for i, team := range project.Teams.Nodes {
					if i > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", team.Key)
				}
				fmt.Println()
			}
			if project.Issues != nil {
				fmt.Printf("Recent Issues: %d shown\n", len(project.Issues.Nodes))
			}
			fmt.Printf("Created: %s\n", project.CreatedAt.Format("2006-01-02"))
			fmt.Printf("Updated: %s\n", project.UpdatedAt.Format("2006-01-02"))
			fmt.Printf("URL: %s\n", project.URL)
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
					color.New(color.FgBlue, color.Underline).Sprint(project.URL))
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
