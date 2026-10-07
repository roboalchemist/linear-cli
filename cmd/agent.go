package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage Linear agent sessions and skills",
	Long: `Inspect Linear agent sessions and workspace agent skills.

Examples:
  linear-cli agent session list
  linear-cli agent session list --limit 10
  linear-cli agent session get SESSION-ID
  linear-cli agent skill list`,
}

var agentSessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"sessions"},
	Short:   "Manage agent sessions",
	Long:    `List and inspect Linear agent sessions.`,
}

var agentSessionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List agent sessions",
	Long:    `List agent sessions in the workspace.`,
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

		sessions, err := client.GetAgentSessions(context.Background(), limit, "", "updatedAt", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list agent sessions: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(sessions.Nodes)
			return
		}

		if len(sessions.Nodes) == 0 {
			output.Info("No agent sessions found", plaintext, jsonOut)
			return
		}

		headers := []string{"Name/ID", "State", "Issue", "Created"}
		rows := make([][]string, len(sessions.Nodes))
		for i, s := range sessions.Nodes {
			name := s.Summary
			if name == "" {
				name = s.SlugID
			}
			if name == "" {
				name = s.ID
			}

			issue := "-"
			if s.Issue != nil {
				issue = s.Issue.Identifier
			}

			state := s.Status
			switch s.Status {
			case "active":
				state = color.New(color.FgGreen).Sprint(s.Status)
			case "complete":
				state = color.New(color.FgBlue).Sprint(s.Status)
			case "error", "stale":
				state = color.New(color.FgRed).Sprint(s.Status)
			case "awaitingInput", "pending", "stopping":
				state = color.New(color.FgYellow).Sprint(s.Status)
			}

			rows[i] = []string{
				truncateString(name, 40),
				state,
				issue,
				s.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		if sessions.PageInfo.HasNextPage && !plaintext {
			fmt.Printf("\n%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

var agentSessionGetCmd = &cobra.Command{
	Use:     "get SESSION-ID",
	Aliases: []string{"show"},
	Short:   "Get agent session details",
	Long:    `Get detailed information about a specific agent session.`,
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
		session, err := client.GetAgentSession(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get agent session: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(session)
			return
		}

		if plaintext {
			fmt.Printf("# Agent Session %s\n", session.ID)
			if session.Summary != "" {
				fmt.Printf("Name: %s\n", session.Summary)
			}
			fmt.Printf("State: %s\n", session.Status)
			if session.Issue != nil {
				fmt.Printf("Issue: %s - %s\n", session.Issue.Identifier, session.Issue.Title)
			}
			if session.AppUser != nil {
				fmt.Printf("App User: %s\n", session.AppUser.Name)
			}
			if session.Creator != nil {
				fmt.Printf("Creator: %s\n", session.Creator.Name)
			}
			if session.StartedAt != nil {
				fmt.Printf("Started: %s\n", session.StartedAt.Format("2006-01-02 15:04:05"))
			}
			if session.EndedAt != nil {
				fmt.Printf("Ended: %s\n", session.EndedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("Created: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05"))
			if session.URL != nil && *session.URL != "" {
				fmt.Printf("URL: %s\n", *session.URL)
			}
			return
		}

		fmt.Printf("\n%s Agent Session\n",
			color.New(color.FgCyan, color.Bold).Sprint("🤖"))
		if session.Summary != "" {
			fmt.Printf("   Name: %s\n", color.New(color.FgWhite, color.Bold).Sprint(session.Summary))
		}
		fmt.Printf("   State: %s\n", session.Status)
		if session.Issue != nil {
			fmt.Printf("   Issue: %s %s\n",
				color.New(color.FgCyan).Sprint(session.Issue.Identifier),
				session.Issue.Title)
		}
		if session.AppUser != nil {
			fmt.Printf("   App User: %s\n", session.AppUser.Name)
		}
		if session.Creator != nil {
			fmt.Printf("   Creator: %s\n", session.Creator.Name)
		}
		fmt.Printf("   Created: %s | Updated: %s\n",
			session.CreatedAt.Format("2006-01-02 15:04:05"),
			session.UpdatedAt.Format("2006-01-02 15:04:05"))
		if session.URL != nil && *session.URL != "" {
			fmt.Printf("   URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(*session.URL))
		}
	},
}

var agentSkillCmd = &cobra.Command{
	Use:     "skill",
	Aliases: []string{"skills"},
	Short:   "Manage agent skills",
	Long:    `List workspace agent skills.`,
}

var agentSkillListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List agent skills",
	Long:    `List all agent skills in the workspace.`,
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

		skills, err := client.GetAgentSkills(context.Background(), limit, "", "createdAt")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list agent skills: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(skills.Nodes)
			return
		}

		if len(skills.Nodes) == 0 {
			output.Info("No agent skills found", plaintext, jsonOut)
			return
		}

		headers := []string{"Name", "Slug", "Shared", "Created"}
		rows := make([][]string, len(skills.Nodes))
		for i, s := range skills.Nodes {
			shared := "no"
			if s.Shared {
				shared = color.New(color.FgGreen).Sprint("yes")
			}
			rows[i] = []string{
				truncateString(s.Title, 40),
				s.SlugID,
				shared,
				s.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		if skills.PageInfo.HasNextPage && !plaintext {
			fmt.Printf("\n%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentSessionCmd)
	agentCmd.AddCommand(agentSkillCmd)
	agentSessionCmd.AddCommand(agentSessionListCmd)
	agentSessionCmd.AddCommand(agentSessionGetCmd)
	agentSkillCmd.AddCommand(agentSkillListCmd)

	agentSessionListCmd.Flags().IntP("limit", "l", 25, "Maximum number of agent sessions to return")
	agentSkillListCmd.Flags().IntP("limit", "l", 25, "Maximum number of agent skills to return")
}
