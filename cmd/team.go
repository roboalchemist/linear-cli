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
	Long: `Manage Linear teams including listing teams, viewing team details, and listing team members.

Examples:
  linear-cli team list              # List all teams
  linear-cli team get ENG           # Get team details
  linear-cli team members ENG       # List team members`,
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

	// List command flags
	teamListCmd.Flags().IntP("limit", "l", 50, "Maximum number of teams to return")
	teamListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
}
