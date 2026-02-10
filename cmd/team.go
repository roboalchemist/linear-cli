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

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage Linear teams",
	Long: `Manage Linear teams including creating, updating, deleting teams, viewing team details, and listing team members.

Examples:
  linear-cli team list              # List all teams
  linear-cli team get ENG           # Get team details
  linear-cli team members ENG       # List team members
  linear-cli team create --name "Engineering" --key ENG
  linear-cli team update ENG --description "Updated"
  linear-cli team delete ENG`,
}

var teamListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List teams",
	Long:    `List all teams in your Linear workspace.`,
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

		// Get teams
		teams, err := client.GetTeams(context.Background(), limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list teams: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(teams.Nodes)
		} else if plaintext {
			fmt.Println("Key\tName\tDescription\tPrivate\tCycles\tTriage\tTimezone\tIssues")
			for _, team := range teams.Nodes {
				description := team.Description
				if len(description) > 50 {
					description = description[:47] + "..."
				}
				cyclesStr := "No"
				if team.CyclesEnabled {
					cyclesStr = "Yes"
				}
				triageStr := "No"
				if team.TriageEnabled {
					triageStr = "Yes"
				}
				fmt.Printf("%s\t%s\t%s\t%v\t%s\t%s\t%s\t%d\n",
					team.Key,
					team.Name,
					description,
					team.Private,
					cyclesStr,
					triageStr,
					team.Timezone,
					team.IssueCount,
				)
			}
		} else {
			// Table output
			headers := []string{"Key", "Name", "Description", "Private", "Cycles", "Triage", "Issues"}
			rows := [][]string{}

			for _, team := range teams.Nodes {
				description := team.Description
				if len(description) > 30 {
					description = description[:27] + "..."
				}

				privateStr := ""
				if team.Private {
					privateStr = color.New(color.FgYellow).Sprint("🔒")
				} else {
					privateStr = color.New(color.FgGreen).Sprint("○")
				}

				cyclesStr := color.New(color.FgRed).Sprint("○")
				if team.CyclesEnabled {
					cyclesStr = color.New(color.FgGreen).Sprint("●")
				}

				triageStr := color.New(color.FgRed).Sprint("○")
				if team.TriageEnabled {
					triageStr = color.New(color.FgGreen).Sprint("●")
				}

				rows = append(rows, []string{
					color.New(color.FgCyan, color.Bold).Sprint(team.Key),
					team.Name,
					description,
					privateStr,
					cyclesStr,
					triageStr,
					fmt.Sprintf("%d", team.IssueCount),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			if !plaintext && !jsonOut {
				fmt.Printf("\n%s %d teams\n",
					color.New(color.FgGreen).Sprint("✓"),
					len(teams.Nodes))
			}
		}
	},
}

var teamGetCmd = &cobra.Command{
	Use:     "get TEAM-KEY",
	Aliases: []string{"show"},
	Short:   "Get team details",
	Long:    `Get detailed information about a specific team.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		teamKey := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get team details
		team, err := client.GetTeam(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get team: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(team)
		} else if plaintext {
			fmt.Printf("Key: %s\n", team.Key)
			fmt.Printf("Name: %s\n", team.Name)
			if team.DisplayName != "" && team.DisplayName != team.Name {
				fmt.Printf("Display Name: %s\n", team.DisplayName)
			}
			if team.Description != "" {
				fmt.Printf("Description: %s\n", team.Description)
			}
			fmt.Printf("Private: %v\n", team.Private)
			fmt.Printf("Issue Count: %d\n", team.IssueCount)
			if team.Timezone != "" {
				fmt.Printf("Timezone: %s\n", team.Timezone)
			}
			if team.Color != "" {
				fmt.Printf("Color: %s\n", team.Color)
			}
			if team.Parent != nil {
				fmt.Printf("Parent Team: %s (%s)\n", team.Parent.Name, team.Parent.Key)
			}
			// Cycle settings
			fmt.Printf("Cycles Enabled: %v\n", team.CyclesEnabled)
			if team.CyclesEnabled {
				fmt.Printf("  Start Day: %d (0=Sun)\n", team.CycleStartDay)
				fmt.Printf("  Duration: %d weeks\n", team.CycleDuration)
				fmt.Printf("  Cooldown: %d weeks\n", team.CycleCooldownTime)
				fmt.Printf("  Upcoming Cycles: %d\n", team.UpcomingCycleCount)
				if team.ActiveCycle != nil {
					fmt.Printf("  Active Cycle: %s (#%d)\n", team.ActiveCycle.Name, team.ActiveCycle.Number)
				}
			}
			// Triage settings
			fmt.Printf("Triage Enabled: %v\n", team.TriageEnabled)
			if team.TriageEnabled && team.TriageIssueState != nil {
				fmt.Printf("  Triage State: %s\n", team.TriageIssueState.Name)
			}
			// Issue estimation
			fmt.Printf("Issue Estimation Type: %s\n", team.IssueEstimationType)
			// Default states
			if team.DefaultIssueState != nil {
				fmt.Printf("Default Issue State: %s\n", team.DefaultIssueState.Name)
			}
			// Auto-archive/close
			fmt.Printf("Auto-Archive Period: %.0f months\n", team.AutoArchivePeriod)
			if team.AutoClosePeriod != nil && *team.AutoClosePeriod > 0 {
				fmt.Printf("Auto-Close Period: %.0f months\n", *team.AutoClosePeriod)
			}
			// AI settings
			if team.AiThreadSummariesEnabled || team.AiDiscussionSummariesEnabled {
				fmt.Printf("AI Thread Summaries: %v\n", team.AiThreadSummariesEnabled)
				fmt.Printf("AI Discussion Summaries: %v\n", team.AiDiscussionSummariesEnabled)
			}
			// Templates
			if team.DefaultTemplateForMembers != nil {
				fmt.Printf("Default Template (Members): %s\n", team.DefaultTemplateForMembers.Name)
			}
			if team.DefaultTemplateForNonMembers != nil {
				fmt.Printf("Default Template (Non-Members): %s\n", team.DefaultTemplateForNonMembers.Name)
			}
			if team.DefaultProjectTemplate != nil {
				fmt.Printf("Default Project Template: %s\n", team.DefaultProjectTemplate.Name)
			}
		} else {
			// Formatted output
			fmt.Println()
			icon := "👥"
			if team.Icon != nil && *team.Icon != "" {
				icon = *team.Icon
			}
			fmt.Printf("%s %s %s (%s)\n",
				color.New(color.FgCyan, color.Bold).Sprint(icon),
				color.New(color.FgCyan, color.Bold).Sprint("Team:"),
				team.Name,
				color.New(color.FgCyan).Sprint(team.Key))
			fmt.Println(strings.Repeat("─", 50))

			if team.Description != "" {
				fmt.Printf("\n%s\n%s\n",
					color.New(color.Bold).Sprint("Description:"),
					team.Description)
			}

			// Basic info section
			fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("Basic Info"))
			privateStr := color.New(color.FgGreen).Sprint("No")
			if team.Private {
				privateStr = color.New(color.FgYellow).Sprint("🔒 Yes")
			}
			fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Private:"), privateStr)
			fmt.Printf("  %s %d\n", color.New(color.Bold).Sprint("Total Issues:"), team.IssueCount)
			if team.Timezone != "" {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Timezone:"), team.Timezone)
			}
			if team.Color != "" {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Color:"), team.Color)
			}
			if team.Parent != nil {
				fmt.Printf("  %s %s (%s)\n", color.New(color.Bold).Sprint("Parent Team:"), team.Parent.Name, team.Parent.Key)
			}

			// Cycles section
			fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("Cycles"))
			cycleStatus := color.New(color.FgRed).Sprint("Disabled")
			if team.CyclesEnabled {
				cycleStatus = color.New(color.FgGreen).Sprint("Enabled")
			}
			fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Status:"), cycleStatus)
			if team.CyclesEnabled {
				days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
				startDay := "Sunday"
				if team.CycleStartDay >= 0 && team.CycleStartDay < 7 {
					startDay = days[team.CycleStartDay]
				}
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Start Day:"), startDay)
				fmt.Printf("  %s %d week(s)\n", color.New(color.Bold).Sprint("Duration:"), team.CycleDuration)
				if team.CycleCooldownTime > 0 {
					fmt.Printf("  %s %d week(s)\n", color.New(color.Bold).Sprint("Cooldown:"), team.CycleCooldownTime)
				}
				if team.ActiveCycle != nil {
					fmt.Printf("  %s %s (#%d)\n", color.New(color.Bold).Sprint("Active Cycle:"),
						color.New(color.FgCyan).Sprint(team.ActiveCycle.Name), team.ActiveCycle.Number)
				}
			}

			// Triage section
			fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("Triage"))
			triageStatus := color.New(color.FgRed).Sprint("Disabled")
			if team.TriageEnabled {
				triageStatus = color.New(color.FgGreen).Sprint("Enabled")
			}
			fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Status:"), triageStatus)
			if team.TriageEnabled && team.TriageIssueState != nil {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Triage State:"), team.TriageIssueState.Name)
			}

			// Estimation section
			fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("Issue Estimation"))
			fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Type:"), team.IssueEstimationType)
			if team.DefaultIssueEstimate > 0 {
				fmt.Printf("  %s %.0f\n", color.New(color.Bold).Sprint("Default Estimate:"), team.DefaultIssueEstimate)
			}

			// Defaults section
			fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("Defaults"))
			if team.DefaultIssueState != nil {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Issue State:"), team.DefaultIssueState.Name)
			}
			if team.DefaultTemplateForMembers != nil {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Template (Members):"), team.DefaultTemplateForMembers.Name)
			}
			if team.DefaultTemplateForNonMembers != nil {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Template (Non-Members):"), team.DefaultTemplateForNonMembers.Name)
			}
			if team.DefaultProjectTemplate != nil {
				fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Project Template:"), team.DefaultProjectTemplate.Name)
			}

			// Automation section
			fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("Automation"))
			fmt.Printf("  %s %.0f months\n", color.New(color.Bold).Sprint("Auto-Archive Period:"), team.AutoArchivePeriod)
			if team.AutoClosePeriod != nil && *team.AutoClosePeriod > 0 {
				fmt.Printf("  %s %.0f months\n", color.New(color.Bold).Sprint("Auto-Close Period:"), *team.AutoClosePeriod)
			}

			// AI section
			if team.AiThreadSummariesEnabled || team.AiDiscussionSummariesEnabled {
				fmt.Printf("\n%s\n", color.New(color.Bold, color.FgYellow).Sprint("AI Features"))
				if team.AiThreadSummariesEnabled {
					fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Thread Summaries:"), color.New(color.FgGreen).Sprint("Enabled"))
				}
				if team.AiDiscussionSummariesEnabled {
					fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("Discussion Summaries:"), color.New(color.FgGreen).Sprint("Enabled"))
				}
			}

			fmt.Println()
		}
	},
}

var teamMembersCmd = &cobra.Command{
	Use:   "members TEAM-KEY",
	Short: "List team members",
	Long:  `List all members of a specific team.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		teamKey := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get team members
		members, err := client.GetTeamMembers(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get team members: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(members.Nodes)
		} else if plaintext {
			fmt.Println("Name\tEmail\tRole\tActive")
			for _, member := range members.Nodes {
				role := "Member"
				if member.Admin {
					role = "Admin"
				}
				fmt.Printf("%s\t%s\t%s\t%v\n",
					member.Name,
					member.Email,
					role,
					member.Active,
				)
			}
		} else {
			// Table output
			headers := []string{"Name", "Email", "Role", "Status"}
			rows := [][]string{}

			for _, member := range members.Nodes {
				role := "Member"
				roleColor := color.New(color.FgWhite)
				if member.Admin {
					role = "Admin"
					roleColor = color.New(color.FgYellow)
				}
				if member.IsMe {
					role = role + " (You)"
					roleColor = color.New(color.FgCyan, color.Bold)
				}

				status := color.New(color.FgGreen).Sprint("✓ Active")
				if !member.Active {
					status = color.New(color.FgRed).Sprint("✗ Inactive")
				}

				rows = append(rows, []string{
					member.Name,
					color.New(color.FgCyan).Sprint(member.Email),
					roleColor.Sprint(role),
					status,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			if !plaintext && !jsonOut {
				fmt.Printf("\n%s %d members in team %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					len(members.Nodes),
					color.New(color.FgCyan).Sprint(teamKey))
			}
		}
	},
}

var teamCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new team",
	Long: `Create a new team in your Linear workspace.

Examples:
  linear-cli team create --name "Engineering"
  linear-cli team create --name "Design" --key DES --description "Design team"
  linear-cli team create --name "Mobile" --color "#4285F4" --private
  linear-cli team create --name "Backend" --copy-settings-from ENG`,
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

		if cmd.Flags().Changed("key") {
			key, _ := cmd.Flags().GetString("key")
			input["key"] = key
		}
		if cmd.Flags().Changed("description") {
			desc, _ := cmd.Flags().GetString("description")
			input["description"] = desc
		}
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}
		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			input["icon"] = icon
		}
		if cmd.Flags().Changed("private") {
			private, _ := cmd.Flags().GetBool("private")
			input["private"] = private
		}
		if cmd.Flags().Changed("cycles-enabled") {
			cycles, _ := cmd.Flags().GetBool("cycles-enabled")
			input["cyclesEnabled"] = cycles
		}
		if cmd.Flags().Changed("triage-enabled") {
			triage, _ := cmd.Flags().GetBool("triage-enabled")
			input["triageEnabled"] = triage
		}
		if cmd.Flags().Changed("timezone") {
			tz, _ := cmd.Flags().GetString("timezone")
			input["timezone"] = tz
		}
		// Cycle settings
		if cmd.Flags().Changed("cycle-start-day") {
			v, _ := cmd.Flags().GetFloat64("cycle-start-day")
			input["cycleStartDay"] = v
		}
		if cmd.Flags().Changed("cycle-duration") {
			v, _ := cmd.Flags().GetInt("cycle-duration")
			input["cycleDuration"] = v
		}
		if cmd.Flags().Changed("cycle-cooldown") {
			v, _ := cmd.Flags().GetInt("cycle-cooldown")
			input["cycleCooldownTime"] = v
		}
		if cmd.Flags().Changed("cycle-auto-assign-started") {
			v, _ := cmd.Flags().GetBool("cycle-auto-assign-started")
			input["cycleIssueAutoAssignStarted"] = v
		}
		if cmd.Flags().Changed("cycle-auto-assign-completed") {
			v, _ := cmd.Flags().GetBool("cycle-auto-assign-completed")
			input["cycleIssueAutoAssignCompleted"] = v
		}
		if cmd.Flags().Changed("cycle-lock-to-active") {
			v, _ := cmd.Flags().GetBool("cycle-lock-to-active")
			input["cycleLockToActive"] = v
		}
		if cmd.Flags().Changed("upcoming-cycle-count") {
			v, _ := cmd.Flags().GetFloat64("upcoming-cycle-count")
			input["upcomingCycleCount"] = v
		}
		// Triage settings
		if cmd.Flags().Changed("require-priority-to-leave-triage") {
			v, _ := cmd.Flags().GetBool("require-priority-to-leave-triage")
			input["requirePriorityToLeaveTriage"] = v
		}
		// Issue estimation settings
		if cmd.Flags().Changed("inherit-issue-estimation") {
			v, _ := cmd.Flags().GetBool("inherit-issue-estimation")
			input["inheritIssueEstimation"] = v
		}
		if cmd.Flags().Changed("issue-estimation-type") {
			v, _ := cmd.Flags().GetString("issue-estimation-type")
			input["issueEstimationType"] = v
		}
		if cmd.Flags().Changed("issue-estimation-allow-zero") {
			v, _ := cmd.Flags().GetBool("issue-estimation-allow-zero")
			input["issueEstimationAllowZero"] = v
		}
		if cmd.Flags().Changed("issue-estimation-extended") {
			v, _ := cmd.Flags().GetBool("issue-estimation-extended")
			input["issueEstimationExtended"] = v
		}
		if cmd.Flags().Changed("default-issue-estimate") {
			v, _ := cmd.Flags().GetFloat64("default-issue-estimate")
			input["defaultIssueEstimate"] = v
		}
		if cmd.Flags().Changed("set-issue-sort-order-on-state-change") {
			v, _ := cmd.Flags().GetString("set-issue-sort-order-on-state-change")
			input["setIssueSortOrderOnStateChange"] = v
		}
		// Workflow settings
		if cmd.Flags().Changed("group-issue-history") {
			v, _ := cmd.Flags().GetBool("group-issue-history")
			input["groupIssueHistory"] = v
		}
		if cmd.Flags().Changed("inherit-workflow-statuses") {
			v, _ := cmd.Flags().GetBool("inherit-workflow-statuses")
			input["inheritWorkflowStatuses"] = v
		}
		// Template settings
		if cmd.Flags().Changed("default-template-for-members") {
			v, _ := cmd.Flags().GetString("default-template-for-members")
			input["defaultTemplateForMembersId"] = v
		}
		if cmd.Flags().Changed("default-template-for-non-members") {
			v, _ := cmd.Flags().GetString("default-template-for-non-members")
			input["defaultTemplateForNonMembersId"] = v
		}
		if cmd.Flags().Changed("default-project-template") {
			v, _ := cmd.Flags().GetString("default-project-template")
			input["defaultProjectTemplateId"] = v
		}
		// Auto-close/archive settings
		if cmd.Flags().Changed("auto-close-period") {
			v, _ := cmd.Flags().GetFloat64("auto-close-period")
			input["autoClosePeriod"] = v
		}
		if cmd.Flags().Changed("auto-close-state") {
			v, _ := cmd.Flags().GetString("auto-close-state")
			input["autoCloseStateId"] = v
		}
		if cmd.Flags().Changed("auto-archive-period") {
			v, _ := cmd.Flags().GetFloat64("auto-archive-period")
			input["autoArchivePeriod"] = v
		}
		if cmd.Flags().Changed("marked-as-duplicate-state") {
			v, _ := cmd.Flags().GetString("marked-as-duplicate-state")
			input["markedAsDuplicateWorkflowStateId"] = v
		}
		// Team hierarchy
		if cmd.Flags().Changed("parent") {
			v, _ := cmd.Flags().GetString("parent")
			input["parentId"] = v
		}

		copyFrom, _ := cmd.Flags().GetString("copy-settings-from")

		team, err := client.CreateTeam(context.Background(), input, copyFrom)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create team: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(team)
		} else if plaintext {
			fmt.Printf("Created team: %s (%s)\n", team.Name, team.Key)
		} else {
			output.Success(fmt.Sprintf("Created team %s (%s)",
				color.New(color.FgWhite, color.Bold).Sprint(team.Name),
				color.New(color.FgCyan).Sprint(team.Key)), plaintext, jsonOut)
		}
	},
}

var teamUpdateCmd = &cobra.Command{
	Use:     "update TEAM-KEY",
	Aliases: []string{"edit"},
	Short:   "Update a team",
	Long: `Update a team's settings.

Examples:
  linear-cli team update ENG --name "Engineering Team"
  linear-cli team update ENG --description "Updated description"
  linear-cli team update ENG --cycles-enabled --triage-enabled
  linear-cli team update ENG --color "#FF5733"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		teamKey := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// First, get the team ID from the key
		team, err := client.GetTeam(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find team: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			input["name"] = name
		}
		if cmd.Flags().Changed("key") {
			key, _ := cmd.Flags().GetString("key")
			input["key"] = key
		}
		if cmd.Flags().Changed("description") {
			desc, _ := cmd.Flags().GetString("description")
			input["description"] = desc
		}
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}
		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			input["icon"] = icon
		}
		if cmd.Flags().Changed("private") {
			private, _ := cmd.Flags().GetBool("private")
			input["private"] = private
		}
		if cmd.Flags().Changed("cycles-enabled") {
			cycles, _ := cmd.Flags().GetBool("cycles-enabled")
			input["cyclesEnabled"] = cycles
		}
		if cmd.Flags().Changed("triage-enabled") {
			triage, _ := cmd.Flags().GetBool("triage-enabled")
			input["triageEnabled"] = triage
		}
		if cmd.Flags().Changed("timezone") {
			tz, _ := cmd.Flags().GetString("timezone")
			input["timezone"] = tz
		}
		// Cycle settings
		if cmd.Flags().Changed("cycle-start-day") {
			v, _ := cmd.Flags().GetFloat64("cycle-start-day")
			input["cycleStartDay"] = v
		}
		if cmd.Flags().Changed("cycle-duration") {
			v, _ := cmd.Flags().GetInt("cycle-duration")
			input["cycleDuration"] = v
		}
		if cmd.Flags().Changed("cycle-cooldown") {
			v, _ := cmd.Flags().GetInt("cycle-cooldown")
			input["cycleCooldownTime"] = v
		}
		if cmd.Flags().Changed("cycle-auto-assign-started") {
			v, _ := cmd.Flags().GetBool("cycle-auto-assign-started")
			input["cycleIssueAutoAssignStarted"] = v
		}
		if cmd.Flags().Changed("cycle-auto-assign-completed") {
			v, _ := cmd.Flags().GetBool("cycle-auto-assign-completed")
			input["cycleIssueAutoAssignCompleted"] = v
		}
		if cmd.Flags().Changed("cycle-lock-to-active") {
			v, _ := cmd.Flags().GetBool("cycle-lock-to-active")
			input["cycleLockToActive"] = v
		}
		if cmd.Flags().Changed("upcoming-cycle-count") {
			v, _ := cmd.Flags().GetFloat64("upcoming-cycle-count")
			input["upcomingCycleCount"] = v
		}
		// Triage settings
		if cmd.Flags().Changed("require-priority-to-leave-triage") {
			v, _ := cmd.Flags().GetBool("require-priority-to-leave-triage")
			input["requirePriorityToLeaveTriage"] = v
		}
		// Issue estimation settings
		if cmd.Flags().Changed("inherit-issue-estimation") {
			v, _ := cmd.Flags().GetBool("inherit-issue-estimation")
			input["inheritIssueEstimation"] = v
		}
		if cmd.Flags().Changed("issue-estimation-type") {
			v, _ := cmd.Flags().GetString("issue-estimation-type")
			input["issueEstimationType"] = v
		}
		if cmd.Flags().Changed("issue-estimation-allow-zero") {
			v, _ := cmd.Flags().GetBool("issue-estimation-allow-zero")
			input["issueEstimationAllowZero"] = v
		}
		if cmd.Flags().Changed("issue-estimation-extended") {
			v, _ := cmd.Flags().GetBool("issue-estimation-extended")
			input["issueEstimationExtended"] = v
		}
		if cmd.Flags().Changed("default-issue-estimate") {
			v, _ := cmd.Flags().GetFloat64("default-issue-estimate")
			input["defaultIssueEstimate"] = v
		}
		if cmd.Flags().Changed("set-issue-sort-order-on-state-change") {
			v, _ := cmd.Flags().GetString("set-issue-sort-order-on-state-change")
			input["setIssueSortOrderOnStateChange"] = v
		}
		// Workflow settings
		if cmd.Flags().Changed("group-issue-history") {
			v, _ := cmd.Flags().GetBool("group-issue-history")
			input["groupIssueHistory"] = v
		}
		if cmd.Flags().Changed("inherit-workflow-statuses") {
			v, _ := cmd.Flags().GetBool("inherit-workflow-statuses")
			input["inheritWorkflowStatuses"] = v
		}
		// Template settings
		if cmd.Flags().Changed("default-template-for-members") {
			v, _ := cmd.Flags().GetString("default-template-for-members")
			input["defaultTemplateForMembersId"] = v
		}
		if cmd.Flags().Changed("default-template-for-non-members") {
			v, _ := cmd.Flags().GetString("default-template-for-non-members")
			input["defaultTemplateForNonMembersId"] = v
		}
		if cmd.Flags().Changed("default-project-template") {
			v, _ := cmd.Flags().GetString("default-project-template")
			input["defaultProjectTemplateId"] = v
		}
		// Auto-close/archive settings
		if cmd.Flags().Changed("auto-close-period") {
			v, _ := cmd.Flags().GetFloat64("auto-close-period")
			input["autoClosePeriod"] = v
		}
		if cmd.Flags().Changed("auto-close-state") {
			v, _ := cmd.Flags().GetString("auto-close-state")
			input["autoCloseStateId"] = v
		}
		if cmd.Flags().Changed("auto-archive-period") {
			v, _ := cmd.Flags().GetFloat64("auto-archive-period")
			input["autoArchivePeriod"] = v
		}
		if cmd.Flags().Changed("auto-close-parent-issues") {
			v, _ := cmd.Flags().GetBool("auto-close-parent-issues")
			input["autoCloseParentIssues"] = v
		}
		if cmd.Flags().Changed("auto-close-child-issues") {
			v, _ := cmd.Flags().GetBool("auto-close-child-issues")
			input["autoCloseChildIssues"] = v
		}
		if cmd.Flags().Changed("marked-as-duplicate-state") {
			v, _ := cmd.Flags().GetString("marked-as-duplicate-state")
			input["markedAsDuplicateWorkflowStateId"] = v
		}
		// Default issue state
		if cmd.Flags().Changed("default-issue-state") {
			v, _ := cmd.Flags().GetString("default-issue-state")
			input["defaultIssueStateId"] = v
		}
		// Team hierarchy
		if cmd.Flags().Changed("parent") {
			v, _ := cmd.Flags().GetString("parent")
			input["parentId"] = v
		}
		// Access/membership settings
		if cmd.Flags().Changed("join-by-default") {
			v, _ := cmd.Flags().GetBool("join-by-default")
			input["joinByDefault"] = v
		}
		if cmd.Flags().Changed("all-members-can-join") {
			v, _ := cmd.Flags().GetBool("all-members-can-join")
			input["allMembersCanJoin"] = v
		}
		if cmd.Flags().Changed("scim-managed") {
			v, _ := cmd.Flags().GetBool("scim-managed")
			input["scimManaged"] = v
		}
		// AI settings
		if cmd.Flags().Changed("ai-thread-summaries") {
			v, _ := cmd.Flags().GetBool("ai-thread-summaries")
			input["aiThreadSummariesEnabled"] = v
		}
		if cmd.Flags().Changed("ai-discussion-summaries") {
			v, _ := cmd.Flags().GetBool("ai-discussion-summaries")
			input["aiDiscussionSummariesEnabled"] = v
		}
		// Slack settings
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

		if len(input) == 0 {
			output.Error("No fields to update. Use --help to see available options.", plaintext, jsonOut)
			os.Exit(1)
		}

		updatedTeam, err := client.UpdateTeam(context.Background(), team.ID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update team: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(updatedTeam)
		} else {
			output.Success(fmt.Sprintf("Updated team %s (%s)",
				color.New(color.FgWhite, color.Bold).Sprint(updatedTeam.Name),
				color.New(color.FgCyan).Sprint(updatedTeam.Key)), plaintext, jsonOut)
		}
	},
}

var teamDeleteCmd = &cobra.Command{
	Use:     "delete TEAM-KEY",
	Aliases: []string{"rm"},
	Short:   "Delete a team",
	Long: `Delete (retire) a team. Teams are soft-deleted with a 30-day grace period during which they can be recovered.

Note: Teams with active issues cannot be deleted.

Examples:
  linear-cli team delete ENG`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		teamKey := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// First, get the team ID from the key
		team, err := client.GetTeam(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find team: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		err = client.DeleteTeam(context.Background(), team.ID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete team: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Team %s deleted", teamKey),
			})
		} else {
			output.Success(fmt.Sprintf("Deleted team %s (30-day recovery period applies)",
				color.New(color.FgCyan).Sprint(teamKey)), plaintext, jsonOut)
		}
	},
}

var teamStatesCmd = &cobra.Command{
	Use:     "states TEAM-KEY",
	Aliases: []string{"workflows"},
	Short:   "List workflow states for a team",
	Long: `List all workflow states (e.g., Triage, Backlog, Todo, In Progress, Done) for a team.
Useful for discovering valid values for the --state flag on issue commands.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		teamKey := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		states, err := client.GetTeamStates(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get team states: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(states)
		} else if plaintext {
			fmt.Println("Name\tType\tColor")
			for _, s := range states {
				fmt.Printf("%s\t%s\t%s\n", s.Name, s.Type, s.Color)
			}
		} else {
			headers := []string{"Name", "Type", "Color"}
			rows := [][]string{}

			for _, s := range states {
				typeColor := color.New(color.FgWhite)
				switch s.Type {
				case "triage":
					typeColor = color.New(color.FgMagenta)
				case "backlog":
					typeColor = color.New(color.FgCyan)
				case "unstarted":
					typeColor = color.New(color.FgWhite)
				case "started":
					typeColor = color.New(color.FgYellow)
				case "completed":
					typeColor = color.New(color.FgGreen)
				case "canceled":
					typeColor = color.New(color.FgRed)
				}

				rows = append(rows, []string{
					typeColor.Sprint(s.Name),
					typeColor.Sprint(s.Type),
					s.Color,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			if !plaintext && !jsonOut {
				fmt.Printf("\n%s %d workflow states for team %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					len(states),
					color.New(color.FgCyan).Sprint(teamKey))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(teamCmd)
	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamGetCmd)
	teamCmd.AddCommand(teamMembersCmd)
	teamCmd.AddCommand(teamStatesCmd)
	teamCmd.AddCommand(teamCreateCmd)
	teamCmd.AddCommand(teamUpdateCmd)
	teamCmd.AddCommand(teamDeleteCmd)

	// List command flags
	teamListCmd.Flags().IntP("limit", "l", 50, "Maximum number of teams to return")
	teamListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")

	// Create command flags - basic settings
	teamCreateCmd.Flags().StringP("name", "n", "", "Team name (required)")
	teamCreateCmd.Flags().StringP("key", "k", "", "Team identifier key (auto-generated from name if omitted)")
	teamCreateCmd.Flags().StringP("description", "d", "", "Team description")
	teamCreateCmd.Flags().StringP("color", "c", "", "Team color (hex, e.g., #4285F4)")
	teamCreateCmd.Flags().String("icon", "", "Team icon")
	teamCreateCmd.Flags().Bool("private", false, "Make the team private")
	teamCreateCmd.Flags().String("timezone", "", "Team timezone")
	teamCreateCmd.Flags().String("copy-settings-from", "", "Copy settings from another team (team key or ID)")
	teamCreateCmd.Flags().String("parent", "", "Parent team ID for team hierarchy")
	_ = teamCreateCmd.MarkFlagRequired("name")

	// Create command flags - cycle settings
	teamCreateCmd.Flags().Bool("cycles-enabled", false, "Enable cycles for the team")
	teamCreateCmd.Flags().Float64("cycle-start-day", 0, "Day of week cycles start (0=Sunday, 1=Monday, etc.)")
	teamCreateCmd.Flags().Int("cycle-duration", 1, "Cycle duration in weeks")
	teamCreateCmd.Flags().Int("cycle-cooldown", 0, "Cooldown time between cycles in weeks")
	teamCreateCmd.Flags().Bool("cycle-auto-assign-started", false, "Auto-assign issues to active cycle when started")
	teamCreateCmd.Flags().Bool("cycle-auto-assign-completed", false, "Auto-assign issues to active cycle when completed")
	teamCreateCmd.Flags().Bool("cycle-lock-to-active", false, "Only allow issues in the active cycle")
	teamCreateCmd.Flags().Float64("upcoming-cycle-count", 2, "Number of upcoming cycles to create")

	// Create command flags - triage settings
	teamCreateCmd.Flags().Bool("triage-enabled", false, "Enable triage mode for the team")
	teamCreateCmd.Flags().Bool("require-priority-to-leave-triage", false, "Require priority to be set before leaving triage")

	// Create command flags - issue estimation settings
	teamCreateCmd.Flags().Bool("inherit-issue-estimation", false, "Inherit issue estimation settings from parent team")
	teamCreateCmd.Flags().String("issue-estimation-type", "", "Issue estimation type (notUsed, exponential, fibonacci, linear, tShirt)")
	teamCreateCmd.Flags().Bool("issue-estimation-allow-zero", false, "Allow zero estimates")
	teamCreateCmd.Flags().Bool("issue-estimation-extended", false, "Use extended estimation scale")
	teamCreateCmd.Flags().Float64("default-issue-estimate", 0, "Default issue estimate")
	teamCreateCmd.Flags().String("set-issue-sort-order-on-state-change", "", "Set issue sort order on state change")

	// Create command flags - workflow settings
	teamCreateCmd.Flags().Bool("group-issue-history", false, "Group issue history entries")
	teamCreateCmd.Flags().Bool("inherit-workflow-statuses", false, "Inherit workflow statuses from parent team")

	// Create command flags - template settings
	teamCreateCmd.Flags().String("default-template-for-members", "", "Default issue template ID for team members")
	teamCreateCmd.Flags().String("default-template-for-non-members", "", "Default issue template ID for non-members")
	teamCreateCmd.Flags().String("default-project-template", "", "Default project template ID")

	// Create command flags - auto-close/archive settings
	teamCreateCmd.Flags().Float64("auto-close-period", 0, "Auto-close period in months (0 to disable)")
	teamCreateCmd.Flags().String("auto-close-state", "", "Workflow state ID to auto-close issues to")
	teamCreateCmd.Flags().Float64("auto-archive-period", 6, "Auto-archive period in months for completed issues")
	teamCreateCmd.Flags().String("marked-as-duplicate-state", "", "Workflow state ID for issues marked as duplicate")

	// Update command flags - basic settings
	teamUpdateCmd.Flags().StringP("name", "n", "", "New team name")
	teamUpdateCmd.Flags().StringP("key", "k", "", "New team identifier key")
	teamUpdateCmd.Flags().StringP("description", "d", "", "New team description")
	teamUpdateCmd.Flags().StringP("color", "c", "", "New team color (hex)")
	teamUpdateCmd.Flags().String("icon", "", "New team icon")
	teamUpdateCmd.Flags().Bool("private", false, "Set team visibility to private")
	teamUpdateCmd.Flags().String("timezone", "", "New team timezone")
	teamUpdateCmd.Flags().String("parent", "", "Parent team ID for team hierarchy")

	// Update command flags - cycle settings
	teamUpdateCmd.Flags().Bool("cycles-enabled", false, "Enable/disable cycles")
	teamUpdateCmd.Flags().Float64("cycle-start-day", 0, "Day of week cycles start (0=Sunday, 1=Monday, etc.)")
	teamUpdateCmd.Flags().Int("cycle-duration", 1, "Cycle duration in weeks")
	teamUpdateCmd.Flags().Int("cycle-cooldown", 0, "Cooldown time between cycles in weeks")
	teamUpdateCmd.Flags().Bool("cycle-auto-assign-started", false, "Auto-assign issues to active cycle when started")
	teamUpdateCmd.Flags().Bool("cycle-auto-assign-completed", false, "Auto-assign issues to active cycle when completed")
	teamUpdateCmd.Flags().Bool("cycle-lock-to-active", false, "Only allow issues in the active cycle")
	teamUpdateCmd.Flags().Float64("upcoming-cycle-count", 2, "Number of upcoming cycles to create")

	// Update command flags - triage settings
	teamUpdateCmd.Flags().Bool("triage-enabled", false, "Enable/disable triage mode")
	teamUpdateCmd.Flags().Bool("require-priority-to-leave-triage", false, "Require priority to be set before leaving triage")

	// Update command flags - issue estimation settings
	teamUpdateCmd.Flags().Bool("inherit-issue-estimation", false, "Inherit issue estimation settings from parent team")
	teamUpdateCmd.Flags().String("issue-estimation-type", "", "Issue estimation type (notUsed, exponential, fibonacci, linear, tShirt)")
	teamUpdateCmd.Flags().Bool("issue-estimation-allow-zero", false, "Allow zero estimates")
	teamUpdateCmd.Flags().Bool("issue-estimation-extended", false, "Use extended estimation scale")
	teamUpdateCmd.Flags().Float64("default-issue-estimate", 0, "Default issue estimate")
	teamUpdateCmd.Flags().String("set-issue-sort-order-on-state-change", "", "Set issue sort order on state change")

	// Update command flags - workflow settings
	teamUpdateCmd.Flags().Bool("group-issue-history", false, "Group issue history entries")
	teamUpdateCmd.Flags().Bool("inherit-workflow-statuses", false, "Inherit workflow statuses from parent team")
	teamUpdateCmd.Flags().String("default-issue-state", "", "Default workflow state ID for new issues")

	// Update command flags - template settings
	teamUpdateCmd.Flags().String("default-template-for-members", "", "Default issue template ID for team members")
	teamUpdateCmd.Flags().String("default-template-for-non-members", "", "Default issue template ID for non-members")
	teamUpdateCmd.Flags().String("default-project-template", "", "Default project template ID")

	// Update command flags - auto-close/archive settings
	teamUpdateCmd.Flags().Float64("auto-close-period", 0, "Auto-close period in months (0 to disable)")
	teamUpdateCmd.Flags().String("auto-close-state", "", "Workflow state ID to auto-close issues to")
	teamUpdateCmd.Flags().Float64("auto-archive-period", 6, "Auto-archive period in months for completed issues")
	teamUpdateCmd.Flags().Bool("auto-close-parent-issues", false, "Auto-close parent issues when all sub-issues are closed")
	teamUpdateCmd.Flags().Bool("auto-close-child-issues", false, "Auto-close child issues when parent is closed")
	teamUpdateCmd.Flags().String("marked-as-duplicate-state", "", "Workflow state ID for issues marked as duplicate")

	// Update command flags - access/membership settings
	teamUpdateCmd.Flags().Bool("join-by-default", false, "New members auto-join this team")
	teamUpdateCmd.Flags().Bool("all-members-can-join", false, "All organization members can join this team")
	teamUpdateCmd.Flags().Bool("scim-managed", false, "Team membership is managed via SCIM")

	// Update command flags - AI settings
	teamUpdateCmd.Flags().Bool("ai-thread-summaries", false, "Enable AI thread summaries")
	teamUpdateCmd.Flags().Bool("ai-discussion-summaries", false, "Enable AI discussion summaries")

	// Update command flags - Slack settings
	teamUpdateCmd.Flags().Bool("slack-new-issue", false, "Notify Slack on new issues")
	teamUpdateCmd.Flags().Bool("slack-issue-comments", false, "Notify Slack on issue comments")
	teamUpdateCmd.Flags().Bool("slack-issue-statuses", false, "Notify Slack on issue status changes")
}
