package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// milestoneCmd is the parent command: project milestone
var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage project milestones",
	Long: `Manage milestones within a Linear project.

Examples:
  linear-cli project milestone list PROJECT-ID
  linear-cli project milestone get MILESTONE-ID
  linear-cli project milestone create PROJECT-ID --name "Beta Release"
  linear-cli project milestone update MILESTONE-ID --name "GA Release"
  linear-cli project milestone delete MILESTONE-ID`,
}

var milestoneListCmd = &cobra.Command{
	Use:     "list PROJECT-ID",
	Aliases: []string{"ls"},
	Short:   "List milestones for a project",
	Long:    `List all milestones within a specific project.`,
	Args:    cobra.ExactArgs(1),
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

		milestones, err := client.GetProjectMilestones(context.Background(), projectID, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list milestones: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(milestones.Nodes) == 0 {
			output.Info("No milestones found for this project.", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(milestones.Nodes)
			return
		}

		if plaintext {
			fmt.Println("# Milestones")
			for _, ms := range milestones.Nodes {
				fmt.Printf("## %s\n", ms.Name)
				fmt.Printf("- **ID**: %s\n", ms.ID)
				fmt.Printf("- **Status**: %s\n", ms.Status)
				fmt.Printf("- **Progress**: %.0f%%\n", ms.Progress*100)
				if ms.TargetDate != nil {
					fmt.Printf("- **Target Date**: %s\n", *ms.TargetDate)
				}
				if ms.Description != nil && *ms.Description != "" {
					fmt.Printf("- **Description**: %s\n", *ms.Description)
				}
				fmt.Println()
			}
			fmt.Printf("Total: %d milestones\n", len(milestones.Nodes))
			return
		}

		// Table output
		headers := []string{"Name", "Status", "Progress", "Target Date", "ID"}
		rows := [][]string{}

		for _, ms := range milestones.Nodes {
			targetDate := "-"
			if ms.TargetDate != nil {
				targetDate = *ms.TargetDate
			}

			statusColor := milestoneStatusColor(ms.Status)

			rows = append(rows, []string{
				truncateString(ms.Name, 30),
				statusColor.Sprint(ms.Status),
				fmt.Sprintf("%.0f%%", ms.Progress*100),
				targetDate,
				ms.ID,
			})
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d milestones\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(milestones.Nodes))
	},
}

var milestoneGetCmd = &cobra.Command{
	Use:     "get MILESTONE-ID",
	Aliases: []string{"show"},
	Short:   "Get milestone details",
	Long:    `Get detailed information about a specific project milestone.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		ms, err := client.GetProjectMilestone(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get milestone: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(ms)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n\n", ms.Name)

			fmt.Printf("## Details\n")
			fmt.Printf("- **ID**: %s\n", ms.ID)
			fmt.Printf("- **Status**: %s\n", ms.Status)
			fmt.Printf("- **Progress**: %.0f%%\n", ms.Progress*100)
			if ms.TargetDate != nil {
				fmt.Printf("- **Target Date**: %s\n", *ms.TargetDate)
			}
			if ms.Description != nil && *ms.Description != "" {
				fmt.Printf("- **Description**: %s\n", *ms.Description)
			}
			fmt.Printf("- **Sort Order**: %.2f\n", ms.SortOrder)
			fmt.Printf("- **Created**: %s\n", ms.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", ms.UpdatedAt.Format("2006-01-02 15:04:05"))
			if ms.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", ms.ArchivedAt.Format("2006-01-02 15:04:05"))
			}

			if ms.Project != nil {
				fmt.Printf("\n## Project\n")
				fmt.Printf("- **Name**: %s\n", ms.Project.Name)
				fmt.Printf("- **State**: %s\n", ms.Project.State)
			}

			if ms.Issues != nil && len(ms.Issues.Nodes) > 0 {
				fmt.Printf("\n## Issues (%d)\n", len(ms.Issues.Nodes))
				for _, issue := range ms.Issues.Nodes {
					stateStr := "[ ]"
					if issue.State != nil {
						switch issue.State.Type {
						case "completed":
							stateStr = "[x]"
						case "started":
							stateStr = "[~]"
						case "canceled":
							stateStr = "[-]"
						}
					}
					assignee := "Unassigned"
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}
					fmt.Printf("- %s %s: %s (%s)\n", stateStr, issue.Identifier, issue.Title, assignee)
				}
			}
			return
		}

		// Rich display
		fmt.Println()
		fmt.Printf("%s %s\n", color.New(color.FgCyan, color.Bold).Sprint("Milestone:"), ms.Name)
		fmt.Println(strings.Repeat("─", 50))

		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), ms.ID)

		statusColor := milestoneStatusColor(ms.Status)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), statusColor.Sprint(ms.Status))

		progressColor := color.New(color.FgRed)
		if ms.Progress >= 0.75 {
			progressColor = color.New(color.FgGreen)
		} else if ms.Progress >= 0.5 {
			progressColor = color.New(color.FgYellow)
		}
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Progress:"), progressColor.Sprintf("%.0f%%", ms.Progress*100))

		if ms.TargetDate != nil {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Target Date:"), *ms.TargetDate)
		}

		if ms.Description != nil && *ms.Description != "" {
			fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Description:"), *ms.Description)
		}

		if ms.Project != nil {
			fmt.Printf("\n%s %s (%s)\n",
				color.New(color.Bold).Sprint("Project:"),
				ms.Project.Name,
				color.New(color.FgWhite, color.Faint).Sprint(ms.Project.State))
		}

		if ms.Issues != nil && len(ms.Issues.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Issues:"))
			for _, issue := range ms.Issues.Nodes {
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

		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Timeline:"))
		fmt.Printf("  Created: %s\n", ms.CreatedAt.Format("2006-01-02"))
		fmt.Printf("  Updated: %s\n", ms.UpdatedAt.Format("2006-01-02"))
		if ms.ArchivedAt != nil {
			fmt.Printf("  Archived: %s\n", ms.ArchivedAt.Format("2006-01-02"))
		}

		fmt.Println()
	},
}

var milestoneCreateCmd = &cobra.Command{
	Use:     "create PROJECT-ID",
	Aliases: []string{"new"},
	Short:   "Create a new milestone",
	Long: `Create a new milestone within a project.

The description can be provided inline via --description or read from a markdown file via --description-file.
Use --description-file - to read from stdin.

Examples:
  linear-cli project milestone create PROJECT-ID --name "Beta Release"
  linear-cli project milestone create PROJECT-ID --name "Beta Release" --description-file milestone-desc.md`,
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

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			output.Error("Name is required (--name)", plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{
			"name":      name,
			"projectId": projectID,
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

		if cmd.Flags().Changed("target-date") {
			targetDate, _ := cmd.Flags().GetString("target-date")
			input["targetDate"] = targetDate
		}

		ms, err := client.CreateProjectMilestone(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create milestone: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(ms)
		} else if plaintext {
			fmt.Printf("Created milestone: %s (ID: %s)\n", ms.Name, ms.ID)
		} else {
			fmt.Printf("%s Created milestone %s (ID: %s)\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(ms.Name),
				ms.ID)
		}
	},
}

var milestoneUpdateCmd = &cobra.Command{
	Use:     "update MILESTONE-ID",
	Aliases: []string{"edit"},
	Short:   "Update a milestone",
	Long: `Update various fields of a project milestone.

The description can be provided inline via --description or read from a markdown file via --description-file.
Use --description-file - to read from stdin.

Examples:
  linear-cli project milestone update MILESTONE-ID --name "New Name"
  linear-cli project milestone update MILESTONE-ID --description-file updated-desc.md
  linear-cli project milestone update MILESTONE-ID --target-date "2025-12-31"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		input := make(map[string]interface{})

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			input["name"] = name
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

		if cmd.Flags().Changed("target-date") {
			targetDate, _ := cmd.Flags().GetString("target-date")
			if targetDate == "" {
				input["targetDate"] = nil
			} else {
				input["targetDate"] = targetDate
			}
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use flags to specify what to update.", plaintext, jsonOut)
			os.Exit(1)
		}

		ms, err := client.UpdateProjectMilestone(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update milestone: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(ms)
		} else if plaintext {
			fmt.Printf("Updated milestone: %s\n", ms.Name)
			fmt.Printf("ID: %s\n", ms.ID)
			fmt.Printf("Status: %s\n", ms.Status)
			fmt.Printf("Progress: %.0f%%\n", ms.Progress*100)
		} else {
			fmt.Printf("%s Updated milestone %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(ms.Name))
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(ms.ID))
			fmt.Printf("  Status: %s\n", ms.Status)
			fmt.Printf("  Progress: %.0f%%\n", ms.Progress*100)
			if ms.TargetDate != nil {
				fmt.Printf("  Target: %s\n", *ms.TargetDate)
			}
		}
	},
}

var milestoneDeleteCmd = &cobra.Command{
	Use:   "delete MILESTONE-ID",
	Short: "Delete a milestone",
	Long:  `Delete a project milestone.`,
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

		err = client.DeleteProjectMilestone(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete milestone: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "id": args[0], "action": "deleted"})
		} else if plaintext {
			fmt.Printf("Deleted milestone %s\n", args[0])
		} else {
			fmt.Printf("%s Deleted milestone %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Faint).Sprint(args[0]))
		}
	},
}

func milestoneStatusColor(status string) *color.Color {
	switch strings.ToLower(status) {
	case "done":
		return color.New(color.FgGreen)
	case "next":
		return color.New(color.FgBlue)
	case "overdue":
		return color.New(color.FgRed)
	default: // unstarted
		return color.New(color.FgWhite)
	}
}

func init() {
	projectCmd.AddCommand(milestoneCmd)
	milestoneCmd.AddCommand(milestoneListCmd)
	milestoneCmd.AddCommand(milestoneGetCmd)
	milestoneCmd.AddCommand(milestoneCreateCmd)
	milestoneCmd.AddCommand(milestoneUpdateCmd)
	milestoneCmd.AddCommand(milestoneDeleteCmd)

	// List flags
	milestoneListCmd.Flags().IntP("limit", "l", 50, "Maximum number of milestones to return")

	// Create flags
	milestoneCreateCmd.Flags().String("name", "", "Milestone name (required)")
	milestoneCreateCmd.Flags().StringP("description", "d", "", "Milestone description")
	milestoneCreateCmd.Flags().String("description-file", "", "Read description from a markdown file (use - for stdin)")
	milestoneCreateCmd.Flags().String("target-date", "", "Target date (YYYY-MM-DD)")
	_ = milestoneCreateCmd.MarkFlagRequired("name")

	// Update flags
	milestoneUpdateCmd.Flags().String("name", "", "New name")
	milestoneUpdateCmd.Flags().StringP("description", "d", "", "New description")
	milestoneUpdateCmd.Flags().String("description-file", "", "Read description from a markdown file (use - for stdin)")
	milestoneUpdateCmd.Flags().String("target-date", "", "New target date (YYYY-MM-DD, or empty to remove)")
}
