package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dorkitude/linear-cli/pkg/api"
	"github.com/dorkitude/linear-cli/pkg/auth"
	"github.com/dorkitude/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Manage Linear cycles (sprints)",
	Long: `Manage Linear cycles including listing current, past, and upcoming cycles.

Examples:
  linear-cli cycle list --team ROB             # List cycles for a team
  linear-cli cycle list --team ROB --active    # Show only the active cycle
  linear-cli cycle get CYCLE-ID                # Get cycle details with issues`,
}

var cycleListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List cycles",
	Long:    `List cycles, optionally filtered by team.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")
		teamKey, _ := cmd.Flags().GetString("team")
		activeOnly, _ := cmd.Flags().GetBool("active")

		filter := map[string]interface{}{}
		if teamKey != "" {
			filter["team"] = map[string]interface{}{
				"key": map[string]interface{}{"eq": teamKey},
			}
		}
		if activeOnly {
			now := time.Now().Format(time.RFC3339)
			filter["startsAt"] = map[string]interface{}{"lte": now}
			filter["endsAt"] = map[string]interface{}{"gte": now}
		}

		cycles, err := client.GetCycles(context.Background(), filter, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list cycles: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(cycles.Nodes)
			return
		}

		if len(cycles.Nodes) == 0 {
			if plaintext {
				fmt.Println("No cycles found")
			} else {
				fmt.Printf("\n%s No cycles found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Cycles")
			fmt.Println("Number\tName\tTeam\tStarts\tEnds\tProgress")
			for _, c := range cycles.Nodes {
				teamName := ""
				if c.Team != nil {
					teamName = c.Team.Key
				}
				fmt.Printf("%d\t%s\t%s\t%s\t%s\t%.0f%%\n",
					c.Number, c.Name, teamName,
					formatDateShort(c.StartsAt), formatDateShort(c.EndsAt),
					c.Progress*100)
			}
		} else {
			headers := []string{"#", "Name", "Team", "Starts", "Ends", "Progress"}
			rows := [][]string{}
			now := time.Now()

			for _, c := range cycles.Nodes {
				teamName := ""
				if c.Team != nil {
					teamName = c.Team.Key
				}

				progressStr := fmt.Sprintf("%.0f%%", c.Progress*100)

				// Determine if active
				startsAt, _ := time.Parse(time.RFC3339, c.StartsAt)
				endsAt, _ := time.Parse(time.RFC3339, c.EndsAt)
				nameStr := c.Name
				if now.After(startsAt) && now.Before(endsAt) {
					nameStr = color.New(color.FgGreen, color.Bold).Sprint(c.Name + " (active)")
				} else if now.Before(startsAt) {
					nameStr = color.New(color.FgCyan).Sprint(c.Name + " (upcoming)")
				} else {
					nameStr = color.New(color.FgWhite, color.Faint).Sprint(c.Name)
				}

				rows = append(rows, []string{
					fmt.Sprintf("%d", c.Number),
					nameStr,
					teamName,
					formatDateShort(c.StartsAt),
					formatDateShort(c.EndsAt),
					progressStr,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
	},
}

var cycleGetCmd = &cobra.Command{
	Use:     "get CYCLE-ID",
	Aliases: []string{"show"},
	Short:   "Get cycle details",
	Long:    `Get details for a specific cycle, including its issues.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		cycleID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		cycle, err := client.GetCycle(context.Background(), cycleID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get cycle: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(cycle)
			return
		}

		if plaintext {
			fmt.Printf("# Cycle %d: %s\n", cycle.Number, cycle.Name)
			if cycle.Description != nil {
				fmt.Printf("Description: %s\n", *cycle.Description)
			}
			fmt.Printf("Team: %s\n", cycle.Team.Key)
			fmt.Printf("Starts: %s\n", formatDateShort(cycle.StartsAt))
			fmt.Printf("Ends: %s\n", formatDateShort(cycle.EndsAt))
			fmt.Printf("Progress: %.0f%%\n", cycle.Progress*100)
			if cycle.Issues != nil {
				fmt.Println("\nIssues:")
				for _, issue := range cycle.Issues.Nodes {
					state := ""
					if issue.State != nil {
						state = issue.State.Name
					}
					assignee := "Unassigned"
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}
					fmt.Printf("  %s\t%s\t%s\t%s\n", issue.Identifier, issue.Title, state, assignee)
				}
			}
		} else {
			teamKey := ""
			if cycle.Team != nil {
				teamKey = cycle.Team.Key
			}
			fmt.Printf("\n%s Cycle %d: %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("🔄"),
				cycle.Number,
				color.New(color.FgWhite, color.Bold).Sprint(cycle.Name))
			fmt.Printf("   Team: %s\n", color.New(color.FgCyan).Sprint(teamKey))
			fmt.Printf("   Period: %s → %s\n", formatDateShort(cycle.StartsAt), formatDateShort(cycle.EndsAt))
			fmt.Printf("   Progress: %s\n",
				color.New(color.FgGreen).Sprintf("%.0f%%", cycle.Progress*100))

			if cycle.Issues != nil && len(cycle.Issues.Nodes) > 0 {
				fmt.Printf("\n   %s Issues:\n\n", color.New(color.FgCyan, color.Bold).Sprint("📋"))
				headers := []string{"ID", "Title", "State", "Assignee"}
				rows := [][]string{}
				for _, issue := range cycle.Issues.Nodes {
					state := ""
					if issue.State != nil {
						state = issue.State.Name
					}
					assignee := "Unassigned"
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}
					rows = append(rows, []string{
						color.New(color.FgCyan).Sprint(issue.Identifier),
						issue.Title,
						state,
						assignee,
					})
				}
				output.Table(output.TableData{
					Headers: headers,
					Rows:    rows,
				}, plaintext, jsonOut)
			}
		}
	},
}

// formatDateShort parses an RFC3339 date and returns YYYY-MM-DD
func formatDateShort(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try date-only format
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return dateStr
		}
	}
	return t.Format("2006-01-02")
}

var cycleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new cycle",
	Long: `Create a new cycle (sprint) for a team.

Examples:
  linear-cli cycle create --team TEAM-ID --name "Sprint 1" --starts 2026-02-10 --ends 2026-02-24`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		teamID, _ := cmd.Flags().GetString("team-id")
		name, _ := cmd.Flags().GetString("name")
		starts, _ := cmd.Flags().GetString("starts")
		ends, _ := cmd.Flags().GetString("ends")

		input := map[string]interface{}{
			"teamId":   teamID,
			"startsAt": starts + "T00:00:00.000Z",
			"endsAt":   ends + "T00:00:00.000Z",
		}
		if name != "" {
			input["name"] = name
		}

		cycle, err := client.CreateCycle(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create cycle: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(cycle)
		} else {
			output.Success(fmt.Sprintf("Created cycle %d: %s (%s → %s)",
				cycle.Number, cycle.Name,
				formatDateShort(cycle.StartsAt), formatDateShort(cycle.EndsAt)), plaintext, jsonOut)
		}
	},
}

var cycleUpdateCmd = &cobra.Command{
	Use:     "update CYCLE-ID",
	Aliases: []string{"edit"},
	Short:   "Update a cycle",
	Long: `Update a cycle's name, description, or dates.

Examples:
  linear-cli cycle update CYCLE-ID --name "Sprint 2"
  linear-cli cycle update CYCLE-ID --starts 2026-03-01 --ends 2026-03-15`,
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
		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			n, _ := cmd.Flags().GetString("name")
			input["name"] = n
		}
		if cmd.Flags().Changed("description") {
			d, _ := cmd.Flags().GetString("description")
			input["description"] = d
		}
		if cmd.Flags().Changed("starts") {
			d, _ := cmd.Flags().GetString("starts")
			input["startsAt"] = d + "T00:00:00.000Z"
		}
		if cmd.Flags().Changed("ends") {
			d, _ := cmd.Flags().GetString("ends")
			input["endsAt"] = d + "T00:00:00.000Z"
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --description, --starts, or --ends.", plaintext, jsonOut)
			os.Exit(1)
		}

		cycle, err := client.UpdateCycle(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update cycle: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(cycle)
		} else {
			output.Success(fmt.Sprintf("Updated cycle %d: %s", cycle.Number, cycle.Name), plaintext, jsonOut)
		}
	},
}

var cycleArchiveCmd = &cobra.Command{
	Use:     "archive CYCLE-ID",
	Aliases: []string{"delete", "rm"},
	Short:   "Archive a cycle",
	Long:    `Archive a cycle. Archived cycles are hidden from the default list.`,
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
		err = client.ArchiveCycle(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to archive cycle: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Archived cycle", plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(cycleCmd)
	cycleCmd.AddCommand(cycleListCmd)
	cycleCmd.AddCommand(cycleGetCmd)
	cycleCmd.AddCommand(cycleCreateCmd)
	cycleCmd.AddCommand(cycleUpdateCmd)
	cycleCmd.AddCommand(cycleArchiveCmd)

	// Update flags
	cycleUpdateCmd.Flags().String("name", "", "New cycle name")
	cycleUpdateCmd.Flags().String("description", "", "New description")
	cycleUpdateCmd.Flags().String("starts", "", "New start date YYYY-MM-DD")
	cycleUpdateCmd.Flags().String("ends", "", "New end date YYYY-MM-DD")

	// List flags
	cycleListCmd.Flags().IntP("limit", "l", 25, "Maximum number of cycles to return")
	cycleListCmd.Flags().StringP("team", "t", "", "Filter by team key (e.g., ROB)")
	cycleListCmd.Flags().Bool("active", false, "Show only the active cycle")

	// Create flags
	cycleCreateCmd.Flags().String("team-id", "", "Team ID (required)")
	cycleCreateCmd.Flags().String("name", "", "Cycle name")
	cycleCreateCmd.Flags().String("starts", "", "Start date YYYY-MM-DD (required)")
	cycleCreateCmd.Flags().String("ends", "", "End date YYYY-MM-DD (required)")
	_ = cycleCreateCmd.MarkFlagRequired("team-id")
	_ = cycleCreateCmd.MarkFlagRequired("starts")
	_ = cycleCreateCmd.MarkFlagRequired("ends")
}
