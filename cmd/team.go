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
			fmt.Println("Key\tName\tDescription\tPrivate\tIssues")
			for _, team := range teams.Nodes {
				description := team.Description
				if len(description) > 50 {
					description = description[:47] + "..."
				}
				fmt.Printf("%s\t%s\t%s\t%v\t%d\n",
					team.Key,
					team.Name,
					description,
					team.Private,
					team.IssueCount,
				)
			}
		} else {
			// Table output
			headers := []string{"Key", "Name", "Description", "Private", "Issues"}
			rows := [][]string{}

			for _, team := range teams.Nodes {
				description := team.Description
				if len(description) > 40 {
					description = description[:37] + "..."
				}

				privateStr := ""
				if team.Private {
					privateStr = color.New(color.FgYellow).Sprint("🔒 Yes")
				} else {
					privateStr = color.New(color.FgGreen).Sprint("No")
				}

				rows = append(rows, []string{
					color.New(color.FgCyan, color.Bold).Sprint(team.Key),
					team.Name,
					description,
					privateStr,
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
			if team.Description != "" {
				fmt.Printf("Description: %s\n", team.Description)
			}
			fmt.Printf("Private: %v\n", team.Private)
			fmt.Printf("Issue Count: %d\n", team.IssueCount)
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s (%s)\n",
				color.New(color.FgCyan, color.Bold).Sprint("👥 Team:"),
				team.Name,
				color.New(color.FgCyan).Sprint(team.Key))
			fmt.Println(strings.Repeat("─", 50))

			if team.Description != "" {
				fmt.Printf("\n%s\n%s\n",
					color.New(color.Bold).Sprint("Description:"),
					team.Description)
			}

			privateStr := color.New(color.FgGreen).Sprint("No")
			if team.Private {
				privateStr = color.New(color.FgYellow).Sprint("🔒 Yes")
			}
			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Private:"), privateStr)
			fmt.Printf("%s %d\n", color.New(color.Bold).Sprint("Total Issues:"), team.IssueCount)
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

		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --key, --description, --color, --icon, --private, --cycles-enabled, --triage-enabled, or --timezone.", plaintext, jsonOut)
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

	// Create command flags
	teamCreateCmd.Flags().StringP("name", "n", "", "Team name (required)")
	teamCreateCmd.Flags().StringP("key", "k", "", "Team identifier key (auto-generated from name if omitted)")
	teamCreateCmd.Flags().StringP("description", "d", "", "Team description")
	teamCreateCmd.Flags().StringP("color", "c", "", "Team color (hex, e.g., #4285F4)")
	teamCreateCmd.Flags().String("icon", "", "Team icon")
	teamCreateCmd.Flags().Bool("private", false, "Make the team private")
	teamCreateCmd.Flags().Bool("cycles-enabled", false, "Enable cycles for the team")
	teamCreateCmd.Flags().Bool("triage-enabled", false, "Enable triage mode for the team")
	teamCreateCmd.Flags().String("timezone", "", "Team timezone")
	teamCreateCmd.Flags().String("copy-settings-from", "", "Copy settings from another team (team key or ID)")
	_ = teamCreateCmd.MarkFlagRequired("name")

	// Update command flags
	teamUpdateCmd.Flags().StringP("name", "n", "", "New team name")
	teamUpdateCmd.Flags().StringP("key", "k", "", "New team identifier key")
	teamUpdateCmd.Flags().StringP("description", "d", "", "New team description")
	teamUpdateCmd.Flags().StringP("color", "c", "", "New team color (hex)")
	teamUpdateCmd.Flags().String("icon", "", "New team icon")
	teamUpdateCmd.Flags().Bool("private", false, "Set team visibility to private")
	teamUpdateCmd.Flags().Bool("cycles-enabled", false, "Enable/disable cycles")
	teamUpdateCmd.Flags().Bool("triage-enabled", false, "Enable/disable triage mode")
	teamUpdateCmd.Flags().String("timezone", "", "New team timezone")
}
