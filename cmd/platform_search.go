package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search across Linear entities",
	Long: `Full-text and semantic search across Linear issues and projects.

Examples:
  linear-cli search issues "payment outage"
  linear-cli search projects "onboarding"
  linear-cli search semantic "billing reliability work"`,
}

var searchIssuesCmd = &cobra.Command{
	Use:     "issues [query]",
	Aliases: []string{"issue"},
	Short:   "Search issues by keyword",
	Long: `Perform a full-text search across Linear issues, ranked by relevance.

Examples:
  linear-cli search issues "auth token"
  linear-cli search issues "customer:" --json`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			output.Error("Search query is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")

		issues, err := client.IssueSearch(context.Background(), query, nil, limit, "", "", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to search issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		emptyMsg := fmt.Sprintf("No matches found for %q", query)
		renderIssueCollection(issues, plaintext, jsonOut, emptyMsg, "matches", "# Search Results")
	},
}

var searchProjectsCmd = &cobra.Command{
	Use:     "projects [query]",
	Aliases: []string{"project"},
	Short:   "Search projects by keyword",
	Long: `Perform a full-text search across Linear projects, ranked by relevance.

Examples:
  linear-cli search projects "onboarding"
  linear-cli search projects "infra" --json`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			output.Error("Search query is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")

		results, err := client.SearchProjects(context.Background(), query, limit, "", "", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to search projects: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(results.Nodes)
			return
		}

		if len(results.Nodes) == 0 {
			output.Info(fmt.Sprintf("No matches found for %q", query), plaintext, jsonOut)
			return
		}

		if plaintext {
			fmt.Println("# Project Search Results")
			for _, p := range results.Nodes {
				fmt.Printf("## %s\n", p.Name)
				fmt.Printf("- **ID**: %s\n", p.ID)
				state := projectSearchState(p)
				if state != "" {
					fmt.Printf("- **State**: %s\n", state)
				}
				fmt.Printf("- **Progress**: %.0f%%\n", p.Progress*100)
				if p.Health != "" {
					fmt.Printf("- **Health**: %s\n", p.Health)
				}
				if p.Lead != nil {
					fmt.Printf("- **Lead**: %s\n", p.Lead.Name)
				}
				fmt.Printf("- **Created**: %s\n", p.CreatedAt.Format("2006-01-02"))
				fmt.Printf("- **URL**: %s\n", p.URL)
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d matches\n", len(results.Nodes))
			return
		}

		headers := []string{"Name", "State", "Progress", "Lead", "Team", "Created", "URL"}
		rows := make([][]string, len(results.Nodes))
		for i, p := range results.Nodes {
			lead := "Unassigned"
			if p.Lead != nil {
				lead = p.Lead.Name
			}
			team := ""
			if p.Teams != nil && len(p.Teams.Nodes) > 0 {
				team = p.Teams.Nodes[0].Key
			}
			rows[i] = []string{
				truncateString(p.Name, 40),
				projectSearchState(p),
				fmt.Sprintf("%.0f%%", p.Progress*100),
				lead,
				team,
				p.CreatedAt.Format("2006-01-02"),
				p.URL,
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d matches\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(results.Nodes))

		if results.PageInfo.HasNextPage {
			fmt.Printf("%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

// projectSearchState prefers the non-deprecated ProjectStatus name and falls
// back to the deprecated state string.
func projectSearchState(p api.ProjectSearchResult) string {
	if p.Status != nil && p.Status.Name != "" {
		return p.Status.Name
	}
	return p.State
}

var searchSemanticCmd = &cobra.Command{
	Use:     "semantic [query]",
	Aliases: []string{"sem"},
	Short:   "Semantic (natural language) search",
	Long: `Search issues, projects, initiatives, and documents using natural-language
semantic search.

Examples:
  linear-cli search semantic "billing reliability work"
  linear-cli search semantic "auth" --types issue --types project`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			output.Error("Search query is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")
		types, _ := cmd.Flags().GetStringSlice("types")

		payload, err := client.SemanticSearch(context.Background(), query, types, limit, false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to run semantic search: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(payload.Results)
			return
		}

		if len(payload.Results) == 0 {
			output.Info(fmt.Sprintf("No matches found for %q", query), plaintext, jsonOut)
			return
		}

		headers := []string{"Type", "Title", "ID"}
		rows := make([][]string, len(payload.Results))
		for i, r := range payload.Results {
			var title, id string
			switch {
			case r.Issue != nil:
				title = fmt.Sprintf("%s - %s", r.Issue.Identifier, r.Issue.Title)
				id = r.Issue.Identifier
			case r.Project != nil:
				title = r.Project.Name
				id = r.Project.ID
			case r.Initiative != nil:
				title = r.Initiative.Name
				id = r.Initiative.ID
			case r.Document != nil:
				title = r.Document.Title
				id = r.Document.ID
			default:
				title = "-"
				id = r.ID
			}
			rows[i] = []string{
				color.New(color.FgCyan).Sprint(r.Type),
				truncateString(title, 60),
				id,
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d matches\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(payload.Results))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(searchIssuesCmd)
	searchCmd.AddCommand(searchProjectsCmd)
	searchCmd.AddCommand(searchSemanticCmd)

	searchIssuesCmd.Flags().IntP("limit", "l", 50, "Maximum number of issues to return")
	searchProjectsCmd.Flags().IntP("limit", "l", 50, "Maximum number of projects to return")
	searchSemanticCmd.Flags().IntP("limit", "l", 25, "Maximum number of results to return")
	searchSemanticCmd.Flags().StringSlice("types", nil, "Restrict to result types (issue, project, initiative, document)")
}
